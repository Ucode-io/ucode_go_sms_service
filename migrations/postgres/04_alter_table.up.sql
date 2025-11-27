ALTER TABLE "sms_send" ADD COLUMN IF NOT EXISTS "originator" VARCHAR;

ALTER TABLE "sms_send" ADD COLUMN mailchimp_key TEXT;
