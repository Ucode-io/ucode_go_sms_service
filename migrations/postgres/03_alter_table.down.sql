ALTER TABLE "sms_send" ADD COLUMN IF NOT EXISTS phone_number VARCHAR(15) NOT NULL;
ALTER TABLE "sms_send" DROP COLUMN IF EXISTS recipient;
ALTER TABLE "sms_send" DROP COLUMN IF EXISTS "type";
ALTER TABLE "sms_send" DROP COLUMN IF EXISTS "dev_email";
ALTER TABLE "sms_send" DROP COLUMN IF EXISTS "dev_email_password";