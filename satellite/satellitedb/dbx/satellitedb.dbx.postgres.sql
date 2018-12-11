-- AUTOGENERATED BY gopkg.in/spacemonkeygo/dbx.v1
-- DO NOT EDIT
CREATE TABLE bwagreements (
	signature bytea NOT NULL,
	data bytea NOT NULL,
	created_at timestamp with time zone NOT NULL,
	PRIMARY KEY ( signature )
);
CREATE TABLE irreparabledbs (
	segmentpath bytea NOT NULL,
	segmentdetail bytea NOT NULL,
	pieces_lost_count bigint NOT NULL,
	seg_damaged_unix_sec bigint NOT NULL,
	repair_attempt_count bigint NOT NULL,
	PRIMARY KEY ( segmentpath )
);