CREATE TABLE APPS (
    ID varchar(36) NOT NULL,
    internal_id int not NULL,
    name varchar(256) NOT NULL,
    ifacts_id varchar(36) NOT NULL,
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    UNIQUE KEY app_internal_id(internal_id),
    PRIMARY KEY APPS_ID(ID)
);

CREATE TABLE APP_PROFILES (
    APP_ID varchar(36) NOT NULL,
    only_handle_centrally boolean DEFAULT FALSE,
    handled_centrally_by varchar(80) DEFAULT '',
    excluded_for_external_supplier boolean DEFAULT FALSE,
    software_development_relevant boolean DEFAULT FALSE,
    cloud_only boolean DEFAULT FALSE,
    physical_security_only boolean DEFAULT FALSE,
    personal_security_only boolean DEFAULT FALSE,
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    UNIQUE KEY app_profiles_id(APP_ID),
    FOREIGN KEY fk_profiles_apps (APP_ID) REFERENCES APPS(ID) ON DELETE CASCADE
);

CREATE TABLE APP_CLASSIFICATIONS (
    IFACTS_ID varchar(36) NOT NULL,
    classification_id varchar(36) NOT NULL,
    classification_name varchar(200) NOT NULL,
    level_id varchar(36) NOT NULL,
    level_name varchar(50) NOT NULL,
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    UNIQUE KEY app_classifications_id(IFACTS_ID, classification_id)
);

CREATE TABLE APP_CONTROLS (
    APP_ID varchar(36) NOT NULL,
    CHECK_ID varchar(16) NOT NULL,
    is_done boolean DEFAULT FALSE,
    notes varchar(10000) NOT NULL DEFAULT '',
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY app_controls_pk(APP_ID, CHECK_ID),
    FOREIGN KEY fk_app_controls_apps (APP_ID) REFERENCES APPS(ID) ON DELETE CASCADE,
    FOREIGN KEY fk_app_controls_checks (CHECK_ID) REFERENCES CHECKS(ID) ON DELETE CASCADE
);

CREATE TRIGGER apps_row_update BEFORE UPDATE
ON APPS
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE TRIGGER app_profiles_row_update BEFORE UPDATE
ON APP_PROFILES
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE TRIGGER app_classifications_row_update BEFORE UPDATE
ON APP_CLASSIFICATIONS
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE TRIGGER app_controls_row_update BEFORE UPDATE
ON APP_CONTROLS
FOR EACH ROW
SET NEW.updated_at = CURRENT_TIMESTAMP;

CREATE VIEW V_APPS AS
    select 
        a.ID,
        a.internal_id,
        a.name, 
        a.ifacts_id,
        ap.only_handle_centrally,
        ap.handled_centrally_by,
        ap.excluded_for_external_supplier,
        ap.software_development_relevant,
        ap.cloud_only,
        ap.physical_security_only,
        ap.personal_security_only
    from APPS a 
    left join APP_PROFILES ap on a.ID = ap.APP_ID;

CREATE VIEW V_APPS_CONTROLS AS
    select
        a.ID as app_id,
        c.ID as check_id,
        c.name,
        c.description,
        ac.is_done,
        ac.notes
    from APP_CONTROLS ac
    left join APPS a on a.ID = ac.APP_ID
    left join CHECKS c on c.ID = ac.CHECK_ID;