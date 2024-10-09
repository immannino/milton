INSTALL sqlite;
LOAD sqlite;

ATTACH 'milton.database' (TYPE SQLITE);
USE milton;

SHOW TABLES;