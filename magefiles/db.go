package main

// // DBGenQueries generates type safe go from sql
// func DBGenQueries() (err error) {
// 	mg.Deps(mg.F(Dep, find))
// 	mg.Deps(mg.F(Dep, sqlc))

// 	schemaStrict := filepath.Join(conf.Dir(), "scripts", "db", "schema.strict.sql")
// 	schema := filepath.Join(conf.Dir(), "scripts", "db", "schema.sql")

// 	// TODO sqlc does not support strict keyword?
// 	// https://github.com/kyleconroy/sqlc/issues/1877
// 	b, err := fileutil.ReadAll(schemaStrict)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	b = bytes.ReplaceAll(b, []byte("strict"), []byte(""))
// 	err = fileutil.WriteBytes(schema, b)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	// Generate code
// 	cmd := exec.Command(sqlc, "compile")
// 	err = printCombinedOutput(cmd)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	cmd = exec.Command(sqlc, "generate")
// 	err = printCombinedOutput(cmd)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	// List generated files
// 	cmd = exec.Command(find,
// 		filepath.Join(conf.Dir(), "go", "dbgen"), "-type", "f", "(",
// 		"-name", "db.go",
// 		"-o", "-name", "models.go",
// 		"-o", "-name", "querier.go",
// 		"-o", "-name", "*.sql.go",
// 		")")
// 	err = printCombinedOutput(cmd)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	return nil
// }
