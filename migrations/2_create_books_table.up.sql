CREATE TABLE `books` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `isbn_13` char(13) DEFAULT NULL,
  `isbn_10` char(10) DEFAULT NULL,
  `publication_year` int(4) NOT NULL,
  `publisher_id` int(11) NOT NULL,
  `image_url` text DEFAULT NULL,
  `edition` varchar(45) DEFAULT NULL,
  `list_price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `books-publisher_id-publisher-id_idx` (`publisher_id`),
  KEY `title_idx` (`title`),
  KEY `isbn_13_idx` (`isbn_13`),
  KEY `isbn_10_idx` (`isbn_10`),
  KEY `publication_year` (`publication_year`),
  CONSTRAINT `books-publisher_id-publishers-id` FOREIGN KEY (`publisher_id`) REFERENCES `publishers` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;