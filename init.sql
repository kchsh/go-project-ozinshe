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
    poster_id       text default null
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

create table users
(
    id            serial primary key,
    name          text not null,
    email         text not null unique,
    password_hash text not null
);

-- Seeding data
insert into users (name, email, password_hash)
values ('admin', 'admin@admin.com', '$2y$10$iCCKNv39bVatC7HelfyfGOLWi9cNYP2zmbb59vIraMMXSnzP5Nczq');

insert into movies(title, description, date_of_release, director, rating, is_watched, trailer_url, poster_id)
values ('1+1',
        'Пострадав в результате несчастного случая, богатый аристократ Филипп нанимает в помощники человека, который менее всего подходит для этой работы, – молодого жителя предместья Дрисса, только что освободившегося из тюрьмы. Несмотря на то, что Филипп прикован к инвалидному креслу, Дриссу удается привнести в размеренную жизнь аристократа дух приключений.',
        '2011-01-01',
        'Оливье Накаш',
        0,
        false,
        'https://www.youtube.com/watch?v=m95M-I7Ij0o&ab_channel=%D0%9A%D0%B8%D0%BD%D0%BE%D0%92%D0%B8%D1%85%D1%80%D1%8C',
        '1+1.jpg'),
       ('Интерстеллар ',
        'Когда засуха, пыльные бури и вымирание растений приводят человечество к продовольственному кризису, коллектив исследователей и учёных отправляется сквозь червоточину (которая предположительно соединяет области пространства-времени через большое расстояние) в путешествие, чтобы превзойти прежние ограничения для космических путешествий человека и найти планету с подходящими для человечества условиями.',
        '2014-01-01',
        'Кристофер Нолан',
        0,
        false,
        'https://www.youtube.com/watch?v=6ybBuTETr3U',
        'Interstellar.jpg'),
       ('Побег из Шоушенка',
        'Бухгалтер Энди Дюфрейн обвинён в убийстве собственной жены и её любовника. Оказавшись в тюрьме под названием Шоушенк, он сталкивается с жестокостью и беззаконием, царящими по обе стороны решётки. Каждый, кто попадает в эти стены, становится их рабом до конца жизни. Но Энди, обладающий живым умом и доброй душой, находит подход как к заключённым, так и к охранникам, добиваясь их особого к себе расположения.',
        '1994-01-01',
        'Фрэнк Дарабонт',
        0,
        false,
        'https://www.youtube.com/watch?v=kgAeKpAPOYk&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC',
        'The Shawshank Redemption.jpg'),
       ('Зеленая миля',
        'Пол Эджкомб — начальник блока смертников в тюрьме «Холодная гора», каждый из узников которого однажды проходит «зеленую милю» по пути к месту казни. Пол повидал много заключённых и надзирателей за время работы. Однако гигант Джон Коффи, обвинённый в страшном преступлении, стал одним из самых необычных обитателей блока.',
        '1999-01-01',
        'Фрэнк Дарабонт',
        0,
        false,
        'https://www.youtube.com/watch?v=TODt_q-_4C4&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC',
        'The Green Mile.jpg'),
       ('Бойцовский клуб',
        'Сотрудник страховой компании страдает хронической бессонницей и отчаянно пытается вырваться из мучительно скучной жизни. Однажды в очередной командировке он встречает некоего Тайлера Дёрдена — харизматического торговца мылом с извращенной философией. Тайлер уверен, что самосовершенствование — удел слабых, а единственное, ради чего стоит жить, — саморазрушение.
Проходит немного времени, и вот уже новые друзья лупят друг друга почем зря на стоянке перед баром, и очищающий мордобой доставляет им высшее блаженство. Приобщая других мужчин к простым радостям физической жестокости, они основывают тайный Бойцовский клуб, который начинает пользоваться невероятной популярностью.',
        '1999-01-01',
        'Дэвид Финчер',
        0,
        false,
        'https://www.youtube.com/watch?v=C7-7qQ61QHU&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC',
        'Fight Club.jpg'),
       ('Остров проклятых',
        'Два американских судебных пристава отправляются на один из островов в штате Массачусетс, чтобы расследовать исчезновение пациентки клиники для умалишенных преступников. При проведении расследования им придется столкнуться с паутиной лжи, обрушившимся ураганом и смертельным бунтом обитателей клиники.',
        '2009-01-01',
        'Мартин Скорсезе',
        0,
        false,
        'https://www.youtube.com/watch?v=_l7R9Rz5URw&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC',
        'Shutter Island.jpg'),
       ('Форрест Гамп',
        'Сидя на автобусной остановке, Форрест Гамп — не очень умный, но добрый и открытый парень — рассказывает случайным встречным историю своей необыкновенной жизни.
С самого малолетства парень страдал от заболевания ног, соседские мальчишки дразнили его, но в один прекрасный день Форрест открыл в себе невероятные способности к бегу. Подруга детства Дженни всегда его поддерживала и защищала, но вскоре дороги их разошлись.',
        '1994-01-01',
        'Роберт Земекис',
        0,
        false,
        'https://www.youtube.com/watch?v=otmeAaifX04',
        'Forrest Gump.jpg'),
       ('Унесённые призраками',
        'Тихиро с мамой и папой переезжает в новый дом. Заблудившись по дороге, они оказываются в странном пустынном городе, где их ждет великолепный пир. Родители с жадностью набрасываются на еду и к ужасу девочки превращаются в свиней, став пленниками злой колдуньи Юбабы. Теперь, оказавшись одна среди волшебных существ и загадочных видений, Тихиро должна придумать, как избавить своих родителей от чар коварной старухи.',
        '2001-01-01',
        'Хаяо Миядзаки',
        0,
        false,
        'https://www.youtube.com/watch?v=bgxiTkAlQrw&ab_channel=iVideos',
        'Sen to Chihiro no kamikakushi.jpg'),
       ('Властелин колец: Возвращение короля',
        'Повелитель сил тьмы Саурон направляет свою бесчисленную армию под стены Минас-Тирита, крепости Последней Надежды. Он предвкушает близкую победу, но именно это мешает ему заметить две крохотные фигурки — хоббитов, приближающихся к Роковой Горе, где им предстоит уничтожить Кольцо Всевластья.',
        '2003-01-01',
        'Питер Джексон',
        0,
        false,
        'https://www.youtube.com/watch?v=lxAeV1-KpSA&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC',
        'lord_of_the_rings.jpg'),
       ('Леон',
        'Профессиональный убийца Леон неожиданно для себя самого решает помочь 12-летней соседке Матильде, семью которой убили коррумпированные полицейские.',
        '1994-01-01',
        'Люк Бессон',
        0,
        false,
        'https://www.youtube.com/watch?v=hvya_q8KM80&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC',
        'leon.jpg');

insert into genres(title)
values ('Драма'),
       ('Комедия'),
       ('Фантастика'),
       ('Приключения'),
       ('Фэнтези'),
       ('Криминал'),
       ('Триллер'),
       ('Детектив'),
       ('Мелодрама'),
       ('Аниме'),
       ('Мультфильм'),
       ('Боевик');

insert into movie_genres(movie_id, genre_id)
values (1, 1),
       (1, 2),
       (2, 1),
       (2, 3),
       (2, 4),
       (3, 1),
       (4, 1),
       (4, 5),
       (4, 6),
       (5, 7),
       (5, 1),
       (5, 6),
       (6, 7),
       (6, 8),
       (6, 1),
       (7, 1),
       (7, 2),
       (7, 9),
       (8, 10),
       (8, 11),
       (8, 5),
       (9, 5),
       (9, 4),
       (9, 1),
       (10, 12),
       (10, 7),
       (10, 1);
       
       
   alter sequence genres_id_seq restart with 1;