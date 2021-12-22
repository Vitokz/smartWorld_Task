CREATE TABLE IF NOT EXISTS users (
    id    SERIAL PRIMARY KEY NOT NULL,
    name varchar NOT NULL DEFAULT '',
    login VARCHAR  NOT NULL UNIQUE DEFAULT '',
    password VARCHAR NOT NULL,
    role  VARCHAR NOT NULL DEFAULT '',
    is_blocked bool NOT NULL DEFAULT false,
    created_at int
);

CREATE TABLE IF NOT EXISTS users_jwt_tokens (
    id_user int UNIQUE NOT NULL,
    access_token VARCHAR DEFAULT '',
    refresh_token VARCHAR  DEFAULT '',
    updated_at int,
    CONSTRAINT user_id_to_jwt FOREIGN KEY (id_user) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE books (
     id SERIAL PRIMARY KEY NOT NULL,
     author VARCHAR NOT NULL DEFAULT '',
     name VARCHAR NOT NULL DEFAULT '',
    count_in_library VARCHAR NOT NULL DEFAULT '',
    updated_at int,
    UNIQUE(name,author)
);

INSERT INTO books (author, name, count_in_library,updated_at)
VALUES
    ('Уильям Голдинг','Повелитель мух',8,123),
    ('Джон Толкин','Властелин колец 1 часть',3,1313),
    ('Джон Толкин','Властелин колец 2 часть',2,1234),
    ('Джон Толкин','Властелин колец 3 часть',1,1423),
    ('Джо Аберкромби','Первый закон 1 книга',4,1234),
    ('Джо Аберкромби','Первый закон 2 книга',1,1234),
    ('Джо Аберкромби','Первый закон 3 книга',2,1423),
    ('Артур Конан Дойл','Багровый этюд',2,1234),
    ('Говард Лавкрафт', 'Зов Ктулху', 3,1234);


CREATE TABLE reservation_books (
    id_user int PRIMARY KEY NOT NULL,
    reservation_books int[],
    CONSTRAINT user_id_to_reservation_books FOREIGN KEY (id_user) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE books_reservation_story(
    id_book     int NOT NULL,
    id_user     int NOT NULL,
    reserved_at int,
    returned_at int,
    CONSTRAINT book_id_to_story_reservation FOREIGN KEY (id_book) REFERENCES books (id) ON UPDATE CASCADE ON DELETE CASCADE
);