-- 既存の 'ctf' データベースを削除（存在する場合）
DROP DATABASE IF EXISTS ctf;

-- 新しく 'ctf' データベースを作成
CREATE DATABASE ctf;

-- 'ctf' データベースを使用
USE ctf;

-- 'teams'
CREATE TABLE teams (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    create_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'contests'
CREATE TABLE contests (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    start        DATETIME NOT NULL,
    end          DATETIME NOT NULL, 
    create_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'users'
CREATE TABLE users (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    name         VARCHAR(100) NOT NULL,
    password     VARCHAR(255) NOT NULL,
    create_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'team_users'
CREATE TABLE team_users (
    team_id        INT UNSIGNED NOT NULL,
    user_id        INT UNSIGNED,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'contest_teams'
CREATE TABLE contest_teams (
    contest_id     INT UNSIGNED NOT NULL,
    team_id        INT UNSIGNED NOT NULL,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id, team_id) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'category'
CREATE TABLE category (
    id             INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name           VARCHAR(255) NOT NULL UNIQUE,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'questions'
CREATE TABLE questions (
    id             INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    category_id    INT UNSIGNED  NOT NULL,
    env            VARCHAR(255),
    description    VARCHAR(255),
    vmid           INT NOT NULL,
    answer         VARCHAR(255),
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'points'
CREATE TABLE points (
    id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    team_id         INT UNSIGNED NOT NULL,
    question_id     INT UNSIGNED NOT NULL,
    contest_id      INT UNSIGNED NOT NULL,
    point           INT UNSIGNED NOT NULL,
    insert_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'contest_questions'
CREATE TABLE contest_questions (
    contest_id INT UNSIGNED NOT NULL,
    question_id INT UNSIGNED NOT NULL,
    point INT NOT NULL,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id,question_id) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- 'cloudinit'
CREATE TABLE cloudinit (
    contest_id INT UNSIGNED NOT NULL,
    question_id INT UNSIGNED NOT NULL,
    team_id                 INT UNSIGNED NOT NULL,
    filename                VARCHAR(255) NOT NULL,
    access                  VARCHAR(255),
    vmid                    INT,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    PRIMARY KEY (contest_id,question_id,team_id), 
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- 'roles'
CREATE TABLE roles (
    id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    namespace       VARCHAR(255) NOT NULL,
    contest         VARCHAR(255)  NOT NULL,
    create_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT unique_name_namespace UNIQUE (name, namespace)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'permissions' (
CREATE TABLE permissions (
    id              INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     VARCHAR(255) NOT NULL,
    create_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'permission_roles'
CREATE TABLE role_permissions (
    role_id INT UNSIGNED NOT NULL,
    permission_id INT UNSIGNED NOT NULL,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 'user_roles'
CREATE TABLE user_roles (
    role_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    create_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_date    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO category (name) VALUES
('Crypto'),
('Reverse Engineering'),
('Web'),
('Forensics');