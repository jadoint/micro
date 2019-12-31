-- MySQL dump 10.13  Distrib 8.0.17, for macos10.14 (x86_64)
--
-- Host: 0.0.0.0    Database: blog
-- ------------------------------------------------------
-- Server version	8.0.17

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `blog`
--

-- DROP TABLE IF EXISTS `blog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog` (
  `id_blog` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `id_author` bigint(20) unsigned NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `post` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `word_count` int(10) unsigned NOT NULL DEFAULT '0',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_blog`),
  KEY `id_author` (`id_author`),
  FULLTEXT KEY `title_fulltext` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`%`*/ /*!50003 TRIGGER `blog_AFTER_INSERT` AFTER INSERT ON `blog` FOR EACH ROW BEGIN
	INSERT INTO blog_settings(id_blog)
	VALUES(new.id_blog);
	
	INSERT INTO blog_views(id_blog)
	VALUES(new.id_blog);
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`%`*/ /*!50003 TRIGGER `blog_BEFORE_DELETE` BEFORE DELETE ON `blog` FOR EACH ROW BEGIN
	DELETE FROM blog_settings
	WHERE id_blog = old.id_blog;
	
	DELETE FROM blog_views
	WHERE id_blog = old.id_blog;
	
	DELETE FROM blog_tags
	WHERE id_blog = old.id_blog;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `blog_settings`
--

-- DROP TABLE IF EXISTS `blog_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog_settings` (
  `id_blog` bigint(20) unsigned NOT NULL,
  `is_unlisted` tinyint(1) NOT NULL DEFAULT '0',
  `is_draft` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_blog`),
  KEY `is_unlisted` (`is_unlisted`),
  CONSTRAINT `blog_settings_fk1` FOREIGN KEY (`id_blog`) REFERENCES `blog` (`id_blog`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `blog_tag`
--

-- DROP TABLE IF EXISTS `blog_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog_tag` (
  `id_tag` bigint(20) unsigned NOT NULL,
  `id_blog` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id_tag`,`id_blog`),
  KEY `id_blog` (`id_blog`),
  CONSTRAINT `blog_tag_fk1` FOREIGN KEY (`id_tag`) REFERENCES `tag` (`id_tag`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `blog_tag_ibfk_2` FOREIGN KEY (`id_blog`) REFERENCES `blog` (`id_blog`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Used for fast retrieval of blogs from a single tag. For blogs under more than one tag, use a FULLTEXT search on the blog_tags table.';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `blog_tags`
--

-- DROP TABLE IF EXISTS `blog_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog_tags` (
  `id_blog` bigint(20) unsigned NOT NULL,
  `tags` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_blog`),
  FULLTEXT KEY `tags_fulltext` (`tags`),
  CONSTRAINT `blog_tags_fk1` FOREIGN KEY (`id_blog`) REFERENCES `blog` (`id_blog`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `blog_views`
--

-- DROP TABLE IF EXISTS `blog_views`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog_views` (
  `id_blog` bigint(20) unsigned NOT NULL,
  `views` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_blog`),
  CONSTRAINT `blog_views_fk1` FOREIGN KEY (`id_blog`) REFERENCES `blog` (`id_blog`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tag`
--

-- DROP TABLE IF EXISTS `tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tag` (
  `id_tag` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `frequency` bigint(20) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id_tag`),
  UNIQUE KEY `tag` (`tag`),
  KEY `frequency` (`frequency`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-09-09 16:12:23
