-- CreateEnum
CREATE TYPE "auth_type_t" AS ENUM ('Google', 'Github', 'Email');

-- CreateTable
CREATE TABLE "User" (
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "picture" TEXT NOT NULL,
    "auth_type" "auth_type_t" NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Question" (
    "id" SERIAL NOT NULL,
    "body" TEXT NOT NULL,
    "driver_code" TEXT,
    "email" TEXT NOT NULL,

    CONSTRAINT "Question_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Testcases" (
    "id" SERIAL NOT NULL,
    "input" TEXT NOT NULL,
    "output" TEXT NOT NULL,
    "order" INTEGER NOT NULL,
    "qid" INTEGER NOT NULL,

    CONSTRAINT "Testcases_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Submission" (
    "id" SERIAL NOT NULL,
    "code" TEXT NOT NULL,
    "question_id" INTEGER NOT NULL,
    "message" TEXT NOT NULL,
    "correct" BOOLEAN NOT NULL,
    "language" TEXT NOT NULL,
    "duration" BIGINT NOT NULL,

    CONSTRAINT "Submission_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "Question_body_key" ON "Question"("body");

-- AddForeignKey
ALTER TABLE "Question" ADD CONSTRAINT "Question_email_fkey" FOREIGN KEY ("email") REFERENCES "User"("email") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Testcases" ADD CONSTRAINT "Testcases_qid_fkey" FOREIGN KEY ("qid") REFERENCES "Question"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
/*
  Warnings:

  - Added the required column `email` to the `Submission` table without a default value. This is not possible if the table is not empty.
  - Added the required column `email` to the `Testcases` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Submission" ADD COLUMN     "email" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "Testcases" ADD COLUMN     "email" TEXT NOT NULL;

-- AddForeignKey
ALTER TABLE "Testcases" ADD CONSTRAINT "Testcases_email_fkey" FOREIGN KEY ("email") REFERENCES "User"("email") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Submission" ADD CONSTRAINT "Submission_email_fkey" FOREIGN KEY ("email") REFERENCES "User"("email") ON DELETE RESTRICT ON UPDATE CASCADE;
