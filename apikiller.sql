CREATE DATABASE `apikiller` character set utf8;

use apikiller;

CREATE TABLE `data_item_strs` (
                                  `id` varchar(50) NOT NULL,
                                  `domain` varchar(100) DEFAULT NULL,
                                  `Url` varchar(500) DEFAULT NULL,
                                  `method` varchar(10) DEFAULT NULL,
                                  `https` tinyint(1) DEFAULT NULL,
                                  `source_request` varchar(50) DEFAULT NULL,
                                  `source_response` varchar(50) DEFAULT NULL,
                                  `vuln_type` varchar(100) DEFAULT NULL,
                                  `vuln_request` varchar(500) DEFAULT NULL,
                                  `vuln_response` varchar(500) DEFAULT NULL,
                                  `check_state` tinyint(1) DEFAULT NULL,
                                  `report_time` varchar(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;



CREATE TABLE `http_items` (
                              `id` int(11) NOT NULL AUTO_INCREMENT,
                              `item` varchar(15000) CHARACTER SET utf8 DEFAULT NULL,
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=latin1;