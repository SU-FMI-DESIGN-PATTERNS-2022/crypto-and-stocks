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

const selectAllWhereUserIdAndSymbolSQL = `
SELECT * FROM "orders" WHERE "user_id"=$1 AND "symbol"=$2`

const updateOrdersAfterMergeSQL = `
UPDATE "orders" SET "user_id"=$1 WHERE "user_id"=$2`
