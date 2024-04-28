create table movies
(
    id           serial primary key,
    title        text not null,
    description  text default null,
    release_date date default null,
    director     text default null,
    rating       float,
    trailer_url  text default null,
    poster_url   text default null
);

create table genres
(
    id    serial primary key,
    title text not null
);

create table movie_genres
(
    movie_id int references movies (id),
    genre_id int references genres (id),
    primary key (movie_id, genre_id)
);