// This file is generated by SQLBoiler (https://github.com/vattle/sqlboiler)
// and is meant to be re-generated in place and/or deleted at any time.
// DO NOT EDIT

package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
	"gopkg.in/nullbio/null.v6"
)

// UserAuth is an object representing the database table.
type UserAuth struct {
	UserAuthID int        `boil:"user_auth_id" json:"user_auth_id" toml:"user_auth_id" yaml:"user_auth_id"`
	UserID     int        `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	Method     string     `boil:"method" json:"method" toml:"method" yaml:"method"`
	Value      null.Bytes `boil:"value" json:"value,omitempty" toml:"value" yaml:"value,omitempty"`

	R *userAuthR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userAuthL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserAuthColumns = struct {
	UserAuthID string
	UserID     string
	Method     string
	Value      string
}{
	UserAuthID: "user_auth_id",
	UserID:     "user_id",
	Method:     "method",
	Value:      "value",
}

// userAuthR is where relationships are stored.
type userAuthR struct {
	User *User
}

// userAuthL is where Load methods for each relationship are stored.
type userAuthL struct{}

var (
	userAuthColumns               = []string{"user_auth_id", "user_id", "method", "value"}
	userAuthColumnsWithoutDefault = []string{"user_id", "method", "value"}
	userAuthColumnsWithDefault    = []string{"user_auth_id"}
	userAuthPrimaryKeyColumns     = []string{"user_auth_id"}
)

type (
	// UserAuthSlice is an alias for a slice of pointers to UserAuth.
	// This should generally be used opposed to []UserAuth.
	UserAuthSlice []*UserAuth

	userAuthQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userAuthType                 = reflect.TypeOf(&UserAuth{})
	userAuthMapping              = queries.MakeStructMapping(userAuthType)
	userAuthPrimaryKeyMapping, _ = queries.BindMapping(userAuthType, userAuthMapping, userAuthPrimaryKeyColumns)
	userAuthInsertCacheMut       sync.RWMutex
	userAuthInsertCache          = make(map[string]insertCache)
	userAuthUpdateCacheMut       sync.RWMutex
	userAuthUpdateCache          = make(map[string]updateCache)
	userAuthUpsertCacheMut       sync.RWMutex
	userAuthUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single userAuth record from the query, and panics on error.
func (q userAuthQuery) OneP() *UserAuth {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single userAuth record from the query.
func (q userAuthQuery) One() (*UserAuth, error) {
	o := &UserAuth{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_auths")
	}

	return o, nil
}

// AllP returns all UserAuth records from the query, and panics on error.
func (q userAuthQuery) AllP() UserAuthSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all UserAuth records from the query.
func (q userAuthQuery) All() (UserAuthSlice, error) {
	var o []*UserAuth

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserAuth slice")
	}

	return o, nil
}

// CountP returns the count of all UserAuth records in the query, and panics on error.
func (q userAuthQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all UserAuth records in the query.
func (q userAuthQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_auths rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q userAuthQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q userAuthQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_auths exists")
	}

	return count > 0, nil
}

// UserG pointed to by the foreign key.
func (o *UserAuth) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *UserAuth) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("user_id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "`users`")

	return query
} // LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userAuthL) LoadUser(e boil.Executor, singular bool, maybeUserAuth interface{}) error {
	var slice []*UserAuth
	var object *UserAuth

	count := 1
	if singular {
		object = maybeUserAuth.(*UserAuth)
	} else {
		slice = *maybeUserAuth.(*[]*UserAuth)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userAuthR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userAuthR{}
			}
			args[i] = obj.UserID
		}
	}

	query := fmt.Sprintf(
		"select * from `users` where `user_id` in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}
	defer results.Close()

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		object.R.User = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.UserID {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// SetUserG of the user_auth to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserAuths.
// Uses the global database handle.
func (o *UserAuth) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the user_auth to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserAuths.
// Panics on error.
func (o *UserAuth) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the user_auth to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserAuths.
// Uses the global database handle and panics on error.
func (o *UserAuth) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the user_auth to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserAuths.
func (o *UserAuth) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `user_auths` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"user_id"}),
		strmangle.WhereClause("`", "`", 0, userAuthPrimaryKeyColumns),
	)
	values := []interface{}{related.UserID, o.UserAuthID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.UserID

	if o.R == nil {
		o.R = &userAuthR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			UserAuths: UserAuthSlice{o},
		}
	} else {
		related.R.UserAuths = append(related.R.UserAuths, o)
	}

	return nil
}

// UserAuthsG retrieves all records.
func UserAuthsG(mods ...qm.QueryMod) userAuthQuery {
	return UserAuths(boil.GetDB(), mods...)
}

// UserAuths retrieves all the records using an executor.
func UserAuths(exec boil.Executor, mods ...qm.QueryMod) userAuthQuery {
	mods = append(mods, qm.From("`user_auths`"))
	return userAuthQuery{NewQuery(exec, mods...)}
}

// FindUserAuthG retrieves a single record by ID.
func FindUserAuthG(userAuthID int, selectCols ...string) (*UserAuth, error) {
	return FindUserAuth(boil.GetDB(), userAuthID, selectCols...)
}

// FindUserAuthGP retrieves a single record by ID, and panics on error.
func FindUserAuthGP(userAuthID int, selectCols ...string) *UserAuth {
	retobj, err := FindUserAuth(boil.GetDB(), userAuthID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindUserAuth retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserAuth(exec boil.Executor, userAuthID int, selectCols ...string) (*UserAuth, error) {
	userAuthObj := &UserAuth{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `user_auths` where `user_auth_id`=?", sel,
	)

	q := queries.Raw(exec, query, userAuthID)

	err := q.Bind(userAuthObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_auths")
	}

	return userAuthObj, nil
}

// FindUserAuthP retrieves a single record by ID with an executor, and panics on error.
func FindUserAuthP(exec boil.Executor, userAuthID int, selectCols ...string) *UserAuth {
	retobj, err := FindUserAuth(exec, userAuthID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *UserAuth) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *UserAuth) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *UserAuth) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *UserAuth) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no user_auths provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(userAuthColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	userAuthInsertCacheMut.RLock()
	cache, cached := userAuthInsertCache[key]
	userAuthInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			userAuthColumns,
			userAuthColumnsWithDefault,
			userAuthColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(userAuthType, userAuthMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userAuthType, userAuthMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `user_auths` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `user_auths` () VALUES ()"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `user_auths` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, userAuthPrimaryKeyColumns))
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into user_auths")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.UserAuthID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == userAuthMapping["UserAuthID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.UserAuthID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for user_auths")
	}

CacheNoHooks:
	if !cached {
		userAuthInsertCacheMut.Lock()
		userAuthInsertCache[key] = cache
		userAuthInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single UserAuth record. See Update for
// whitelist behavior description.
func (o *UserAuth) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single UserAuth record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *UserAuth) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the UserAuth, and panics on error.
// See Update for whitelist behavior description.
func (o *UserAuth) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the UserAuth.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *UserAuth) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	userAuthUpdateCacheMut.RLock()
	cache, cached := userAuthUpdateCache[key]
	userAuthUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			userAuthColumns,
			userAuthPrimaryKeyColumns,
			whitelist,
		)

		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update user_auths, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `user_auths` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, userAuthPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userAuthType, userAuthMapping, append(wl, userAuthPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update user_auths row")
	}

	if !cached {
		userAuthUpdateCacheMut.Lock()
		userAuthUpdateCache[key] = cache
		userAuthUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q userAuthQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q userAuthQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for user_auths")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UserAuthSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o UserAuthSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o UserAuthSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserAuthSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userAuthPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `user_auths` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userAuthPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in userAuth slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *UserAuth) UpsertG(updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *UserAuth) UpsertGP(updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *UserAuth) UpsertP(exec boil.Executor, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *UserAuth) Upsert(exec boil.Executor, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no user_auths provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(userAuthColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userAuthUpsertCacheMut.RLock()
	cache, cached := userAuthUpsertCache[key]
	userAuthUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			userAuthColumns,
			userAuthColumnsWithDefault,
			userAuthColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			userAuthColumns,
			userAuthPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert user_auths, could not build update column list")
		}

		cache.query = queries.BuildUpsertQueryMySQL(dialect, "user_auths", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `user_auths` WHERE `user_auth_id`=?",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)

		cache.valueMapping, err = queries.BindMapping(userAuthType, userAuthMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userAuthType, userAuthMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for user_auths")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.UserAuthID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == userAuthMapping["UserAuthID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.UserAuthID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for user_auths")
	}

CacheNoHooks:
	if !cached {
		userAuthUpsertCacheMut.Lock()
		userAuthUpsertCache[key] = cache
		userAuthUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single UserAuth record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *UserAuth) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single UserAuth record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *UserAuth) DeleteG() error {
	if o == nil {
		return errors.New("models: no UserAuth provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single UserAuth record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *UserAuth) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single UserAuth record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserAuth) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no UserAuth provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userAuthPrimaryKeyMapping)
	sql := "DELETE FROM `user_auths` WHERE `user_auth_id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from user_auths")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q userAuthQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q userAuthQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no userAuthQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from user_auths")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o UserAuthSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o UserAuthSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no UserAuth slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o UserAuthSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserAuthSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no UserAuth slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userAuthPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `user_auths` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userAuthPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from userAuth slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *UserAuth) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *UserAuth) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *UserAuth) ReloadG() error {
	if o == nil {
		return errors.New("models: no UserAuth provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserAuth) Reload(exec boil.Executor) error {
	ret, err := FindUserAuth(exec, o.UserAuthID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *UserAuthSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *UserAuthSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserAuthSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty UserAuthSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserAuthSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	userAuths := UserAuthSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userAuthPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `user_auths`.* FROM `user_auths` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userAuthPrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&userAuths)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserAuthSlice")
	}

	*o = userAuths

	return nil
}

// UserAuthExists checks if the UserAuth row exists.
func UserAuthExists(exec boil.Executor, userAuthID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `user_auths` where `user_auth_id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, userAuthID)
	}

	row := exec.QueryRow(sql, userAuthID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_auths exists")
	}

	return exists, nil
}

// UserAuthExistsG checks if the UserAuth row exists.
func UserAuthExistsG(userAuthID int) (bool, error) {
	return UserAuthExists(boil.GetDB(), userAuthID)
}

// UserAuthExistsGP checks if the UserAuth row exists. Panics on error.
func UserAuthExistsGP(userAuthID int) bool {
	e, err := UserAuthExists(boil.GetDB(), userAuthID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// UserAuthExistsP checks if the UserAuth row exists. Panics on error.
func UserAuthExistsP(exec boil.Executor, userAuthID int) bool {
	e, err := UserAuthExists(exec, userAuthID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
