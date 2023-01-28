package user_repository

const createUserSQL = `
INSERT INTO "users"("userID", "name", "orders", "isBot", "creatorID", "amount")
VALUES($1,$2,$3,$4,$5,$6)`

const createBotSQL = `
INSERT INTO "users"("userID", "name", "orders", "isBot", "creatorID", "amount")
VALUES($1,$2,$3,$4,$5,$6)`

const selectAmountWhereIdSQL = `
SELECT "amount" FROM "users" WHERE "id"=$1`

const selectOrdersWhereIdSQL = `
SELECT "orders" FROM "users" WHERE "id"=$1`

const selectAllWhereCreatorIdSQL = `
SELECT * FROM "users" WHERE "creatorId"=$1 IN (SELECT * FROM "users" WHERE "isBot"=true)`

const updateUserOrdersWhereIdSQL = `
UPDATE "users" SET "orders"=$1 WHERE "id"=$2`

const deleteUserWhereIdSQL = `
DELETE FROM "users" WHERE "Id"=$1`
