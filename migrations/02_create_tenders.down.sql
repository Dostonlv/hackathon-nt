DROP INDEX IF EXISTS idx_tenders_client_id;
DROP INDEX IF EXISTS idx_tenders_status;
DROP TABLE IF EXISTS tenders;
DROP TYPE IF EXISTS tender_status;

DROP INDEX IF EXISTS idx_bids_tender_id;
DROP INDEX IF EXISTS idx_bids_contractor_id;
DROP TABLE IF EXISTS bids;

DROP INDEX IF EXISTS idx_notifications_user_id;
DROP TABLE IF EXISTS notifications;