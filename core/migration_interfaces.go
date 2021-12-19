package core

import (
	"container/list"
	"fmt"
	"sort"
	"strings"
)

type IMigrationRegistry interface {
	GetByName(migrationName string) (IMigration, error)
	AddMigration(migration IMigration)
	GetSortedMigrations() MigrationList
}

type IMigration interface {
	Up(database *ProjectDatabase) error
	Down(database *ProjectDatabase) error
	GetName() string
	GetID() int64
	Deps() []string
}

type IMigrationNode interface {
	IsApplied() bool
	GetMigration() IMigration
	SetItAsRoot()
	IsRoot() bool
	AddChild(node IMigrationNode)
	AddDep(node IMigrationNode)
	GetChildrenCount() int
	GetChildren() *list.List
	GetDeps() *list.List
	TraverseDeps(migrationList []int64, depList MigrationDepList) MigrationDepList
	TraverseChildren(migrationList []int64) []int64
	IsDummy() bool
	Downgrade(database *ProjectDatabase) error
	Apply(database *ProjectDatabase) error
}

type IMigrationTree interface {
	GetRoot() IMigrationNode
	SetRoot(root IMigrationNode)
	GetNodeByMigrationID(ID int64) (IMigrationNode, error)
	AddNode(node IMigrationNode) error
	TreeBuilt()
	IsTreeBuilt() bool
}

func GetBluePrintNameFromMigrationName(migrationName string) string {
	return strings.Split(migrationName, ".")[0]
}

type MigrationNode struct {
	Deps     *list.List
	Node     IMigration
	Children *list.List
	applied  bool
	isRoot   bool
	dummy    bool
}

func (n MigrationNode) IsDummy() bool {
	return n.dummy
}

func (n MigrationNode) Apply(database *ProjectDatabase) error {
	res := n.Node.Up(database)
	if res == nil {
		n.applied = true
	}
	return res
}

func (n MigrationNode) Downgrade(database *ProjectDatabase) error {
	res := n.Node.Down(database)
	if res == nil {
		n.applied = false
	}
	return res
}

func (n MigrationNode) GetMigration() IMigration {
	return n.Node
}

func (n MigrationNode) GetChildren() *list.List {
	return n.Children
}

func (n MigrationNode) GetDeps() *list.List {
	return n.Deps
}

func (n MigrationNode) GetChildrenCount() int {
	return n.Children.Len()
}

func (n MigrationNode) IsApplied() bool {
	return n.applied
}

func (n MigrationNode) SetItAsRoot() {
	n.isRoot = true
}

func (n MigrationNode) IsRoot() bool {
	return n.isRoot
}

func (n MigrationNode) AddChild(node IMigrationNode) {
	n.Children.PushBack(node)
}

func (n MigrationNode) AddDep(node IMigrationNode) {
	n.Deps.PushBack(node)
}

func (n MigrationNode) TraverseDeps(migrationList []int64, depList MigrationDepList) MigrationDepList {
	for l := n.GetDeps().Front(); l != nil; l = l.Next() {
		migration := l.Value.(IMigrationNode)
		if migration.IsDummy() {
			continue
		}
		migrationName := l.Value.(IMigrationNode).GetMigration().GetID()
		if !ContainsInt64(migrationList, migrationName) && !ContainsInt64(depList, migrationName) {
			depList = append(depList, l.Value.(IMigrationNode).GetMigration().GetID())
			depList = l.Value.(IMigrationNode).TraverseDeps(migrationList, depList)
		}
	}
	return depList
}

func (n MigrationNode) TraverseChildren(migrationList []int64) []int64 {
	for l := n.GetChildren().Front(); l != nil; l = l.Next() {
		migration := l.Value.(IMigrationNode)
		if migration.IsDummy() {
			continue
		}
		migrationName := l.Value.(IMigrationNode).GetMigration().GetID()
		if !ContainsInt64(migrationList, migrationName) {
			migrationList = append(migrationList, l.Value.(IMigrationNode).GetMigration().GetID())
			migrationDepList := n.TraverseDeps(migrationList, make(MigrationDepList, 0))
			sort.Slice(migrationDepList, func(i int, j int) bool {
				return i > j
			})
			for _, m := range migrationDepList {
				migrationList = append(migrationList, m)
			}
			migrationList = n.TraverseChildren(migrationList)
		}
	}
	return migrationList
}

func NewMigrationNode(dep IMigrationNode, node IMigration, child IMigrationNode) IMigrationNode {
	depsList := list.New()
	if dep != nil {
		depsList.PushBack(dep)
	}
	childrenList := list.New()
	if child != nil {
		childrenList.PushBack(child)
	}
	return &MigrationNode{
		Deps:     depsList,
		Node:     node,
		Children: childrenList,
		applied:  false,
		dummy:    false,
		isRoot:   false,
	}
}

func NewMigrationRootNode() IMigrationNode {
	return &MigrationNode{
		Deps:     list.New(),
		Node:     nil,
		Children: list.New(),
		applied:  false,
		dummy:    true,
		isRoot:   true,
	}
}

type MigrationTree struct {
	Root      IMigrationNode
	nodes     map[int64]IMigrationNode
	treeBuilt bool
}

func (t MigrationTree) TreeBuilt() {
	t.treeBuilt = true
}

func (t MigrationTree) IsTreeBuilt() bool {
	return t.treeBuilt
}

func (t MigrationTree) GetNodeByMigrationID(migrationID int64) (IMigrationNode, error) {
	node, ok := t.nodes[migrationID]
	if ok {
		return node, nil
	}
	return nil, fmt.Errorf("no node with name %d has been found", migrationID)
}

func (t MigrationTree) AddNode(node IMigrationNode) error {
	_, ok := t.nodes[node.GetMigration().GetID()]
	if ok {
		// return fmt.Errorf("Migration with name %s has been added to tree before", node.GetMigration().GetName())
		return nil
	}
	t.nodes[node.GetMigration().GetID()] = node
	return nil
}

func (t MigrationTree) GetRoot() IMigrationNode {
	return t.Root
}

func (t MigrationTree) SetRoot(root IMigrationNode) {
	root.SetItAsRoot()
	t.Root = root
}

type MigrationList []IMigration

func (m MigrationList) Len() int { return len(m) }
func (m MigrationList) Less(i, j int) bool {
	return m[i].GetID() < m[j].GetID()
}
func (m MigrationList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

type MigrationDepList []int64

func (m MigrationDepList) Len() int { return len(m) }
func (m MigrationDepList) Less(i, j int) bool {
	return i < j
}
func (m MigrationDepList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

type MigrationRegistry struct {
	migrations map[string]IMigration
}

func (r MigrationRegistry) AddMigration(migration IMigration) {
	r.migrations[migration.GetName()] = migration
}

func (r MigrationRegistry) GetByName(migrationName string) (IMigration, error) {
	migration, ok := r.migrations[migrationName]
	if ok {
		return migration, nil
	}
	return nil, fmt.Errorf("No migration with name %s exists", migrationName)
}

func (r MigrationRegistry) GetSortedMigrations() MigrationList {
	sortedMigrations := make(MigrationList, len(r.migrations))
	i := 0
	for _, migration := range r.migrations {
		sortedMigrations[i] = migration
		i++
	}
	sort.Slice(sortedMigrations, func(i int, j int) bool {
		return sortedMigrations[i].GetID() < sortedMigrations[j].GetID()
	})
	return sortedMigrations
}

func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		migrations: make(map[string]IMigration),
	}
}

func NewMigrationTree() IMigrationTree {
	return &MigrationTree{
		Root:  NewMigrationRootNode(),
		nodes: make(map[int64]IMigrationNode),
	}
}
