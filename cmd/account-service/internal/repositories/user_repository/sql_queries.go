package user_repository

const createUserSQL = `
INSERT INTO "users"("user_id", "name", "is_bot", "creator_id", "amount")
VALUES($1,$2,$3,$4,$5)`

const createBotSQL = `
INSERT INTO "users"("user_id", "name", "is_bot", "creator_id", "amount")
VALUES($1,$2,$3,$4,$5)`

const selectUserWhereIdSQL = `
SELECT * FROM "users" WHERE "id"=$1`

const selectUserAmountWhereIdSQL = `
SELECT "amount" FROM "users" WHERE "id"=$1`

const selectUserWhereUserIdSQL = `
SELECT * FROM "users" WHERE "user_id"=$1`

const updateUserAmountSQL = `
UPDATE "users" SET "amount"=$1 WHERE "id"=$2`

const deleteUserWhereIdSQL = `
DELETE FROM "users" WHERE "id"=$1`
