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

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");
-- CreateTable
CREATE TABLE "Question" (
    "id" SERIAL NOT NULL,
    "body" TEXT NOT NULL,
    "driver_code" TEXT,
    "input_testcases" TEXT NOT NULL,
    "output" TEXT NOT NULL,

    CONSTRAINT "Question_pkey" PRIMARY KEY ("id")
);
/*
  Warnings:

  - A unique constraint covering the columns `[body]` on the table `Question` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `email` to the `Question` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Question" ADD COLUMN     "email" TEXT NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "Question_body_key" ON "Question"("body");

-- AddForeignKey
ALTER TABLE "Question" ADD CONSTRAINT "Question_email_fkey" FOREIGN KEY ("email") REFERENCES "User"("email") ON DELETE RESTRICT ON UPDATE CASCADE;
/*
  Warnings:

  - You are about to drop the column `input_testcases` on the `Question` table. All the data in the column will be lost.
  - You are about to drop the column `output` on the `Question` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "Question" DROP COLUMN "input_testcases",
DROP COLUMN "output";

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

    CONSTRAINT "Submission_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "Testcases" ADD CONSTRAINT "Testcases_qid_fkey" FOREIGN KEY ("qid") REFERENCES "Question"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
-- AlterTable
ALTER TABLE "Testcases" ALTER COLUMN "output" DROP NOT NULL;
/*
  Warnings:

  - Made the column `output` on table `Testcases` required. This step will fail if there are existing NULL values in that column.

*/
-- AlterTable
ALTER TABLE "Testcases" ALTER COLUMN "output" SET NOT NULL;


/*
  Warnings:

  - Added the required column `correct` to the `Submission` table without a default value. This is not possible if the table is not empty.
  - Added the required column `message` to the `Submission` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Submission" ADD COLUMN     "correct" BOOLEAN NOT NULL,
ADD COLUMN     "message" TEXT NOT NULL;


/*
  Warnings:

  - Added the required column `language` to the `Submission` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Submission" ADD COLUMN     "language" TEXT NOT NULL;
/*
  Warnings:

  - Added the required column `duration` to the `Submission` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Submission" ADD COLUMN     "duration" BIGINT NOT NULL;
