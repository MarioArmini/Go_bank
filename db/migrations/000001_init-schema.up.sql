CREATE TABLE "Accounts" (
  "Id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "creationTime" timestamptz NOT NULL DEFAULT 'now()',
  "countryCode" int,
  "interestRate" DECIMAL(5, 2) NOT NULL DEFAULT 0.05
);

CREATE TABLE "Entries" (
  "Id" bigserial PRIMARY KEY,
  "accountId" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "creationTime" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "Transfers" (
  "Id" bigserial PRIMARY KEY,
  "senderId" bigint NOT NULL,
  "recipientId" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "creationTime" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "Accounts" ("owner");

CREATE INDEX ON "Entries" ("accountId");

CREATE INDEX ON "Transfers" ("senderId");

CREATE INDEX ON "Transfers" ("recipientId");

CREATE INDEX ON "Transfers" ("senderId", "recipientId");

ALTER TABLE "Entries" ADD FOREIGN KEY ("accountId") REFERENCES "Accounts" ("Id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("senderId") REFERENCES "Accounts" ("Id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("recipientId") REFERENCES "Accounts" ("Id");