# Migrate command

Use this command to manipulate data migrations for your project.

# Create new migration for your blueprint

Execute following command:
```bash
  ./gomonolith_binary migrate create -m "Message that describes the idea of your data migration" -b "{{SHORT_NAME_OF_YOUR_BLUEPRINT}}"
```
As result of this command execution, you would have a new migration created for your blueprint.  

# Update your database

Execute following command:
```bash
  ./gomonolith_binary migrate up
```
All not applied data migrations will be applied.

# Downgrade your database

Execute following command:
```bash
  ./gomonolith_binary migrate down --to-id {{MIGRATION_ID - each migration has ID (int64)}}
```
All migrations applied after this migration ID would be downgraded

# Determine conflicts in parallel migration branches

Execute following command:
```bash
  ./gomonolith_binary migrate determine-conflicts
```
Whole migration tree would be traversed and all parallel migration branches within one blueprint would be determined.

# Create migrations to merge all conflicted branches

Execute following command:
```bash
  ./gomonolith_binary migrate create --merge
```
This command would create migrations to merge all conflicted branches. **KEEP IN MIND. IT DOESN'T SOLVE CONFLICTS**, but it creates leaf that merges all conflicted branches.  
So, please make sure investigate all changes and make sure you resolve all conflicts. Ideally, if there's no real conflicts, then, this migration would be used to merge all together all parallel branches.
