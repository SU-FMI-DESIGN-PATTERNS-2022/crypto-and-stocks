package user_repository

const createUserSQL = `
INSERT INTO "users"("user_id", "name", "is_bot", "creator_id", "amount")
VALUES($1,$2,$3,$4,$5)`

const createBotSQL = `
INSERT INTO "users"("user_id", "name", "is_bot", "creator_id", "amount")
VALUES($1,$2,$3,$4,$5)`

const selectUserWhereIdSQL = `
SELECT * FROM "users" WHERE "id"=$1`

const selectAllOrdersWhereUserIdSQL = `
SELECT * FROM "orders" WHERE "user_id"=$1`

const selectAllWhereCreatorIdSQL = `
SELECT * FROM "users" WHERE "creator_id"=$1 AND "is_bot"=true`

const updateOrdersAfterMergeSQL = `
UPDATE "orders" SET "user_id"=$1 WHERE "user_id"=$2`

const updateUserAmountSQL = `
UPDATE "users" SET "amount"=$1 WHERE "id"=$2`

const deleteUserWhereIdSQL = `
DELETE FROM "users" WHERE "id"=$1`
