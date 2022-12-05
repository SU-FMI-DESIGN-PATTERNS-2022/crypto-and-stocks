package user_repository

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
