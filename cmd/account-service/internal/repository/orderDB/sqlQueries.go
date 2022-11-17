package orderDB

const insertSQL = `
INSERT INTO "orders"("userID", "type", "symbol", "amount", "price", "date")
VALUES($1,$2,$3,$4,$5,$6)`

const selectAllSQL = `
SELECT * FROM "orders"`

const selectAllWhereSymbolSQL = `
SELECT * FROM "orders" WHERE "symbol"=$1`

const selectAllWhereUserIdSQL = `
SELECT * FROM "orders" WHERE "userId"=$1`
