package core

import (
	"encoding/json"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"sort"
)

type TraverseMigrationResult struct {
	Node  IMigrationNode
	Error error
}

type IBlueprint interface {
	GetName() string
	GetDescription() string
	GetMigrationRegistry() IMigrationRegistry
	InitRouter(app IApp, group *gin.RouterGroup)
	InitApp(app IApp)
}

type IBlueprintRegistry interface {
	Iterate() <-chan IBlueprint
	GetByName(name string) (IBlueprint, error)
	Register(blueprint IBlueprint)
	GetMigrationTree() IMigrationTree
	TraverseMigrations() <-chan *TraverseMigrationResult
	TraverseMigrationsDownTo(downToMigration int64) <-chan *TraverseMigrationResult
	InitializeRouting(app IApp, router *gin.Engine)
	Initialize(app IApp)
	ResetMigrationTree()
	DeRegister(blueprint IBlueprint)
}

type Blueprint struct {
	Name              string
	Description       string
	MigrationRegistry IMigrationRegistry
}

func (b Blueprint) GetName() string {
	return b.Name
}

func (b Blueprint) InitRouter(app IApp, group *gin.RouterGroup) {
	panic(errors.New("has to be redefined in concrete blueprint"))
}

func (b Blueprint) GetDescription() string {
	return b.Description
}

func (b Blueprint) InitApp(app IApp) {

}

func (b Blueprint) GetMigrationRegistry() IMigrationRegistry {
	return b.MigrationRegistry
}

type BlueprintRegistry struct {
	RegisteredBlueprints map[string]IBlueprint
	MigrationTree        IMigrationTree
}

func (r BlueprintRegistry) ResetMigrationTree() {
	r.MigrationTree = NewMigrationTree()
}

func (r BlueprintRegistry) Iterate() <-chan IBlueprint {
	chnl := make(chan IBlueprint)
	go func() {
		defer close(chnl)
		for _, blueprint := range r.RegisteredBlueprints {
			chnl <- blueprint
		}
	}()
	return chnl
}

func (r BlueprintRegistry) GetByName(name string) (IBlueprint, error) {
	blueprint, ok := r.RegisteredBlueprints[name]
	var err error
	if !ok {
		err = fmt.Errorf("Couldn't find blueprint with name %s", name)
	}
	return blueprint, err
}

func (r BlueprintRegistry) Register(blueprint IBlueprint) {
	r.RegisteredBlueprints[blueprint.GetName()] = blueprint
}

func (r BlueprintRegistry) DeRegister(blueprint IBlueprint) {
	delete(r.RegisteredBlueprints, blueprint.GetName())
}

func (r BlueprintRegistry) GetMigrationTree() IMigrationTree {
	return r.MigrationTree
}

func (r BlueprintRegistry) buildMigrationTree(chnl chan *TraverseMigrationResult) bool {
	if r.MigrationTree.IsTreeBuilt() {
		return true
	}
	var err error
	_, err = r.GetByName("user")
	if err != nil {
		res := &TraverseMigrationResult{
			Node:  nil,
			Error: err,
		}
		chnl <- res
		return false
	}
	var currentMigration IMigration
	var tmpMigration IMigration
	var previousNode IMigrationNode
	var currentNode IMigrationNode
	var tmpNode IMigrationNode
	var blueprintName string
	var blueprintTmp IBlueprint
	for blueprint := range r.Iterate() {
		numberOfMigrationsWithNoDeps := mapset.NewSet()
		for i, migration := range blueprint.GetMigrationRegistry().GetSortedMigrations() {
			currentMigration = migration
			if i == 0 {
				if numberOfMigrationsWithNoDeps.Cardinality() >= 1 {
					numberOfMigrationsWithNoDeps.Add(currentMigration.GetName())
					res := &TraverseMigrationResult{
						Node:  nil,
						Error: fmt.Errorf("Found two migrations with no deps : %v", numberOfMigrationsWithNoDeps),
					}
					chnl <- res
					return false
				}
				currentNode, err = r.MigrationTree.GetNodeByMigrationID(currentMigration.GetID())
				if currentNode == nil {
					currentNode = NewMigrationNode(
						r.MigrationTree.GetRoot(), currentMigration, nil,
					)
				}
				err = r.MigrationTree.AddNode(currentNode)
				if err != nil {
					res := &TraverseMigrationResult{
						Node:  nil,
						Error: err,
					}
					chnl <- res
					return false
				}
				r.MigrationTree.GetRoot().AddChild(
					currentNode,
				)
				numberOfMigrationsWithNoDeps.Add(currentMigration.GetName())
			}
			if i != 0 {
				currentNode, err = r.MigrationTree.GetNodeByMigrationID(currentMigration.GetID())
				if currentNode == nil {
					currentNode = NewMigrationNode(
						previousNode, currentMigration, nil,
					)
				}
				// previousNode.AddChild(currentNode)
				err = r.MigrationTree.AddNode(currentNode)
				if err != nil {
					res := &TraverseMigrationResult{
						Node:  nil,
						Error: err,
					}
					chnl <- res
					return false
				}
			}
			for _, dep := range currentMigration.Deps() {
				blueprintName = GetBluePrintNameFromMigrationName(dep)
				if blueprintName == blueprint.GetName() {
					tmpMigration, err = blueprint.GetMigrationRegistry().GetByName(dep)
					if err != nil {
						res := &TraverseMigrationResult{
							Node:  nil,
							Error: err,
						}
						chnl <- res
						return false
					}
				} else {
					blueprintTmp, err = r.GetByName(blueprintName)
					if err != nil {
						res := &TraverseMigrationResult{
							Node:  nil,
							Error: err,
						}
						chnl <- res
						return false
					}
					tmpMigration, err = blueprintTmp.GetMigrationRegistry().GetByName(dep)
					if err != nil {
						res := &TraverseMigrationResult{
							Node:  nil,
							Error: err,
						}
						chnl <- res
						return false
					}
				}
				tmpNode, err = r.MigrationTree.GetNodeByMigrationID(tmpMigration.GetID())
				if tmpNode == nil {
					tmpNode = NewMigrationNode(
						previousNode, tmpMigration, nil,
					)
				}
				currentNode.AddDep(tmpNode)
				tmpNode.AddChild(currentNode)
			}
			previousNode = currentNode
		}
	}
	for blueprint := range r.Iterate() {
		numberOfMigrationsWithNoDescendants := mapset.NewSet()
		for _, migration := range blueprint.GetMigrationRegistry().GetSortedMigrations() {
			currentNode, err = r.MigrationTree.GetNodeByMigrationID(migration.GetID())
			if err != nil {

			}
			if currentNode.GetChildrenCount() == 0 {
				numberOfMigrationsWithNoDescendants.Add(migration.GetName())
			}
		}
		if numberOfMigrationsWithNoDescendants.Cardinality() >= 2 {
			res := &TraverseMigrationResult{
				Node:  nil,
				Error: fmt.Errorf("Found two or more migrations with no children from the same blueprint: %v", numberOfMigrationsWithNoDescendants),
			}
			chnl <- res
			return false
		}
	}
	return true
}

func (r BlueprintRegistry) TraverseMigrations() <-chan *TraverseMigrationResult {
	chnl := make(chan *TraverseMigrationResult)
	go func() {
		defer close(chnl)
		wasTreeBuilt := r.buildMigrationTree(chnl)
		if !wasTreeBuilt {
			return
		}
		r.MigrationTree.TreeBuilt()
		applyMigrationsInOrder := make([]int64, 0)
		allBlueprintRoots := r.MigrationTree.GetRoot().GetChildren()
		for l := allBlueprintRoots.Front(); l != nil; l = l.Next() {
			applyMigrationsInOrder = append(applyMigrationsInOrder, l.Value.(IMigrationNode).GetMigration().GetID())
			migrationDepList := l.Value.(IMigrationNode).TraverseDeps(applyMigrationsInOrder, make(MigrationDepList, 0))
			sort.Slice(migrationDepList, func(i int, j int) bool {
				return i > j
			})
			for _, m := range migrationDepList {
				applyMigrationsInOrder = append(applyMigrationsInOrder, m)
			}
			applyMigrationsInOrder = l.Value.(IMigrationNode).TraverseChildren(applyMigrationsInOrder)
		}
		for _, migrationName := range applyMigrationsInOrder {
			node, err := r.MigrationTree.GetNodeByMigrationID(migrationName)
			if err != nil {
				res := &TraverseMigrationResult{
					Node:  nil,
					Error: fmt.Errorf("not found migration node with id : %d", migrationName),
				}
				chnl <- res
				return
			}
			res := &TraverseMigrationResult{
				Node:  node,
				Error: nil,
			}
			chnl <- res
		}
	}()
	return chnl
}

func (r BlueprintRegistry) InitializeRouting(app IApp, router *gin.Engine) {
	for blueprint := range r.Iterate() {
		routergroup := router.Group("/" + blueprint.GetName())
		blueprint.InitRouter(app, routergroup)
	}
	router.GET("/ping/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/localization/", func(ctx *gin.Context) {
		type Context struct {
			AdminContext
		}
		c := &Context{}
		adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
		PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)
		langMap := readLocalization(c.GetLanguage().Code)
		langMapB, _ := json.Marshal(langMap)
		ctx.Header("Content-Type", "application/javascript")
		ctx.String(200, fmt.Sprintf("setLocalization(%s)", string(langMapB)))
	})
	if CurrentConfig.D.Debug {
		router.GET("/get_not_translated/", func(ctx *gin.Context) {
			notTranslated, _ := json.Marshal(NotTranslatedData.D)
			ctx.Header("Content-Type", "application/json")
			ctx.String(200, string(notTranslated))
		})
	}
	router.POST("/testcsrf/", func(c *gin.Context) {
		c.String(200, "csrf token test passed")
	})
	router.POST("/ignorecsrfcheck/", func(c *gin.Context) {
		c.String(200, "csrf token test passed")
	})
}

func (r BlueprintRegistry) Initialize(app IApp) {
	ClearProjectModels()
	ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &ContentType{}, &[]*ContentType{} })
	for blueprint := range r.Iterate() {
		blueprint.InitApp(app)
	}
}

func (r BlueprintRegistry) TraverseMigrationsDownTo(downToMigration int64) <-chan *TraverseMigrationResult {
	// @todo, fix this implementation
	chnl := make(chan *TraverseMigrationResult)
	go func() {
		defer close(chnl)
		wasTreeBuilt := r.buildMigrationTree(chnl)
		if !wasTreeBuilt {
			return
		}
		applyMigrationsInOrder := make([]int64, 0)
		allBlueprintRoots := r.MigrationTree.GetRoot().GetChildren()
		for l := allBlueprintRoots.Front(); l != nil; l = l.Next() {
			applyMigrationsInOrder = append(applyMigrationsInOrder, l.Value.(IMigrationNode).GetMigration().GetID())
			migrationDepList := l.Value.(IMigrationNode).TraverseDeps(applyMigrationsInOrder, make(MigrationDepList, 0))
			sort.Slice(migrationDepList, func(i int, j int) bool {
				return i > j
			})
			for _, m := range migrationDepList {
				applyMigrationsInOrder = append(applyMigrationsInOrder, m)
			}
			applyMigrationsInOrder = l.Value.(IMigrationNode).TraverseChildren(applyMigrationsInOrder)
		}
		downgradeMigrationsInOrder := make(MigrationDepList, 0)
		var foundMigration = false
		for _, migrationName := range applyMigrationsInOrder {
			if downToMigration > 0 && migrationName == downToMigration {
				foundMigration = true
			}
			if downToMigration > 0 && !foundMigration {
				continue
			}
			downgradeMigrationsInOrder = append(downgradeMigrationsInOrder, migrationName)
		}
		sort.Slice(downgradeMigrationsInOrder, func(i int, j int) bool {
			return i > j
		})
		for _, migrationName := range downgradeMigrationsInOrder {
			node, err := r.MigrationTree.GetNodeByMigrationID(migrationName)
			if err != nil {
				res := &TraverseMigrationResult{
					Node:  nil,
					Error: fmt.Errorf("not found migration node with name : %d", migrationName),
				}
				chnl <- res
				return
			}
			res := &TraverseMigrationResult{
				Node:  node,
				Error: nil,
			}
			chnl <- res
		}
	}()
	return chnl
}

func NewBlueprintRegistry() IBlueprintRegistry {
	return &BlueprintRegistry{
		RegisteredBlueprints: make(map[string]IBlueprint),
		MigrationTree:        NewMigrationTree(),
	}
}
