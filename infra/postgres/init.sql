-- DEV ENVIRONMENT ONLY — credentials below are for local development
-- App user: nutrometra / nutrometra_dev

-- Banco da aplicação
CREATE DATABASE nutrometra
    WITH ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE   = 'en_US.utf8';

-- Banco do Zitadel
CREATE DATABASE zitadel
    WITH ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE   = 'en_US.utf8';

-- Usuário dedicado para a aplicação
CREATE USER nutrometra WITH PASSWORD 'nutrometra_dev';
GRANT ALL PRIVILEGES ON DATABASE nutrometra TO nutrometra;
