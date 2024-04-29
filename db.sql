create table movies
(
    id              serial primary key,
    title           text not null,
    description     text default null,
    date_of_release date default null,
    director        text default null,
    rating          int  default 0,
    is_watched      bool default false,
    trailer_url     text default null,
    poster_url      text default null
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

create table watchlist
(
    movie_id int primary key references movies (id),
    added_at timestamp
);