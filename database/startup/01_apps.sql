CREATE TABLE APPS (
    ID varchar(36) NOT NULL,
    internal_id int not NULL,
    name varchar(256) NOT NULL,
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
    APP_ID varchar(36) NOT NULL,
    c int NOT NULL DEFAULT 1,
    i int NOT NULL DEFAULT 1,
    a int NOT NULL DEFAULT 1,
    t int NOT NULL DEFAULT 1,
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    UNIQUE KEY app_classifications_id(APP_ID),
    FOREIGN KEY fk_classifications_apps (APP_ID) REFERENCES APPS(ID) ON DELETE CASCADE
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

CREATE VIEW V_APPS AS
    select 
        a.ID,
        a.internal_id,
        a.name, 
        ap.only_handle_centrally,
        ap.handled_centrally_by,
        ap.excluded_for_external_supplier,
        ap.software_development_relevant,
        ap.cloud_only,
        ap.physical_security_only,
        ap.personal_security_only,
        ac.c,
        ac.i,
        ac.a,
        ac.t
    from APPS a 
    inner join APP_PROFILES ap on a.ID = ap.APP_ID
    inner join APP_CLASSIFICATIONS ac on a.ID = ac.APP_ID;