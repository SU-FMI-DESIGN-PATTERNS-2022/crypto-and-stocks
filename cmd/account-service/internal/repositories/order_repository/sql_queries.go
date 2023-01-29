package order_repository

const insertSQL = `
INSERT INTO "orders"("user_id", "type", "symbol", "amount", "price", "date")
VALUES($1,$2,$3,$4,$5,$6)`

const selectAllSQL = `
SELECT * FROM "orders"`

const selectAllWhereSymbolSQL = `
SELECT * FROM "orders" WHERE "symbol"=$1`

const selectAllWhereUserIdSQL = `
SELECT * FROM "orders" WHERE "user_id"=$1`

const selectUserAmountWhereUserIdSQL = `
SELECT "amount" FROM "users" WHERE "id"=$1`
