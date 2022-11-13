package orderDB

const insertSQL = `
INSERT INTO "orders"("userID", "type", "symbol", "amount", "price", "date")
VALUES($1,$2,$3,$4,$5,$6)`

const selectAllSQL = `
SELECT * FROM "orders"`
