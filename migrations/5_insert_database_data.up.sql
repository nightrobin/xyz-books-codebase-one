INSERT INTO `authors` (`first_name`, `middle_name`, `last_name`) VALUES ('Joel', '', 'Hartse');
INSERT INTO `authors` (`first_name`, `middle_name`, `last_name`) VALUES ('Hannah', 'P.', 'Templer');
INSERT INTO `authors` (`first_name`, `middle_name`, `last_name`) VALUES ('Marguerite', 'Z.', 'Duras');
INSERT INTO `authors` (`first_name`, `last_name`) VALUES ('Kingsley', 'Amis');
INSERT INTO `authors` (`first_name`, `middle_name`, `last_name`) VALUES ('Fannie', 'Peters', 'Flagg');
INSERT INTO `authors` (`first_name`, `middle_name`, `last_name`) VALUES ('Camille', 'Byron', 'Paglia');
INSERT INTO `authors` (`first_name`, `middle_name`, `last_name`) VALUES ('Rainer', 'Steel', 'Rilke');

INSERT INTO `publishers` (`id`, `name`) VALUES ('1', 'Paste Magazine');
INSERT INTO `publishers` (`id`, `name`) VALUES ('2', 'Publishers Weekly');
INSERT INTO `publishers` (`id`, `name`) VALUES ('3', 'Graywolf Press');
INSERT INTO `publishers` (`id`, `name`) VALUES ('4', 'McSweeney\'s');

INSERT INTO `books` (`id`, `title`, `isbn_13`, `isbn_10`, `publication_year`, `publisher_id`, `edition`, `list_price`) VALUES ('1', 'American Elf', '9781891830853', '1891830856', '2004', '1', 'Book 2', '1000.00');
INSERT INTO `books` (`id`, `title`, `isbn_13`, `isbn_10`, `publication_year`, `publisher_id`, `edition`, `list_price`) VALUES ('2', 'Cosmoknights', '9781603094542', '1603094547', '2019', '2', 'Book 1', '2000.00');
INSERT INTO `books` (`id`, `title`, `isbn_13`, `isbn_10`, `publication_year`, `publisher_id`, `list_price`) VALUES ('3', 'Essex County', '9781603090384', '160309038X', '1990', '3', '500.00');
INSERT INTO `books` (`id`, `title`, `isbn_13`, `isbn_10`, `publication_year`, `publisher_id`, `edition`, `list_price`) VALUES ('4', 'Hey, Mister (Vol 1)', '9781891830020', '1891830023', '2000', '3', 'After School Special', '1200.00');
INSERT INTO `books` (`id`, `title`, `isbn_13`, `isbn_10`, `publication_year`, `publisher_id`, `list_price`) VALUES ('5', 'The Underwater Welder', '9781603093989', '1603093982', '2022', '4', '3000.00');

INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('1', '1', '1');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('2', '1', '2');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('3', '1', '3');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('4', '2', '4');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('5', '3', '4');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('6', '4', '2');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('7', '4', '5');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('8', '4', '6');
INSERT INTO `book_authors` (`id`, `book_id`, `author_id`) VALUES ('9', '5', '7');
