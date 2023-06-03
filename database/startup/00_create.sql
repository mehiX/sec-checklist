CREATE TABLE CHECKS (
    ID varchar(16) NOT NULL,
    type varchar(50) NOT NULL,
    name varchar(150) NOT NULL,
    description varchar(4000) DEFAULT '',
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY id(ID)
);

CREATE TABLE CHECKS_CIAT (
    CHECK_ID varchar(16) NOT NULL,
    c varchar(16) NULL DEFAULT '',
    i varchar(16) NULL DEFAULT '',
    a varchar(16) NULL DEFAULT '',
    t varchar(16) NULL DEFAULT '',
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    UNIQUE KEY ciat_check_id(CHECK_ID),
    FOREIGN KEY fk_ciat_checks (CHECK_ID) REFERENCES CHECKS(ID) ON DELETE CASCADE
);

CREATE TABLE FILTERS (
    CHECK_ID varchar(16) NOT NULL,
    only_handle_centrally boolean DEFAULT FALSE,
    handled_centrally_by varchar(80) DEFAULT '',
    excluded_for_external_supplier boolean DEFAULT FALSE,
    software_development_relevant boolean DEFAULT FALSE,
    cloud_only boolean DEFAULT FALSE,
    physical_security_only boolean DEFAULT FALSE,
    personal_security_only boolean DEFAULT FALSE,
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    UNIQUE KEY filters_check_id(CHECK_ID),
    FOREIGN KEY fk_filters_checks (CHECK_ID) REFERENCES CHECKS(ID) ON DELETE CASCADE
);

CREATE TRIGGER checks_row_update BEFORE UPDATE
ON CHECKS
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE TRIGGER checks_ciat_row_update BEFORE UPDATE
ON CHECKS_CIAT
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE TRIGGER filters_row_update BEFORE UPDATE
ON FILTERS
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE VIEW V_CHECKS AS
    select 
        a.ID, a.type, a.name, a.description, b.c, b.i, b.a, b.t, 
        f.only_handle_centrally,
        f.handled_centrally_by,
        f.excluded_for_external_supplier,
        f.software_development_relevant,
        f.cloud_only,
        f.physical_security_only,
        f.personal_security_only
    from CHECKS a 
    inner join CHECKS_CIAT b on a.ID = b.CHECK_ID
    inner join FILTERS f on a.ID = f.CHECK_ID;
    
