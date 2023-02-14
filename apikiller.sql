CREATE DATABASE `apikiller` character set utf8;

use apikiller;

CREATE TABLE `data_item_strs` (
                                  `id` varchar(50) NOT NULL,
                                  `domain` varchar(100) DEFAULT NULL,
                                  `url` varchar(500) DEFAULT NULL,
                                  `method` varchar(20) DEFAULT NULL,
                                  `https` tinyint(1) DEFAULT NULL,
                                  `source_request` varchar(5000) DEFAULT NULL,
                                  `source_response` varchar(5000) DEFAULT NULL,
                                  `vuln_type` varchar(20) DEFAULT NULL,
                                  `vuln_request` varchar(5000) DEFAULT NULL,
                                  `vuln_response` varchar(5000) DEFAULT NULL,
                                  `check_time` varchar(10) DEFAULT NULL,
                                  `check_state` tinyint(1) DEFAULT NULL,
                                  `report_time` varchar(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;