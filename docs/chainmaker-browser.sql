-- MySQL dump 10.13  Distrib 8.3.0, for macos14.2 (x86_64)
--
-- Host: localhost    Database: explorer_format_2
-- ------------------------------------------------------
-- Server version	8.0.33

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `chainmaker_pk_account`
--

DROP TABLE IF EXISTS `chainmaker_pk_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_account` (
  `address` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `addrType` bigint DEFAULT NULL,
  `did` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bns` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`address`),
  KEY `did_index` (`did`),
  KEY `bns_index` (`bns`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_black_transaction`
--

DROP TABLE IF EXISTS `chainmaker_pk_black_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_black_transaction` (
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sender` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `senderOrgId` longtext COLLATE utf8mb4_unicode_ci,
  `blockHeight` bigint DEFAULT NULL,
  `blockHash` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `txType` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `expirationTime` bigint DEFAULT NULL,
  `txIndex` bigint DEFAULT NULL,
  `txStatusCode` longtext COLLATE utf8mb4_unicode_ci,
  `rwSetHash` longtext COLLATE utf8mb4_unicode_ci,
  `contractResultCode` int unsigned DEFAULT NULL,
  `contractResult` mediumblob,
  `contractResultBak` mediumblob,
  `contractMessage` blob,
  `contractMessageBak` blob,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `contractNameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractRuntimeType` longtext COLLATE utf8mb4_unicode_ci,
  `contractType` longtext COLLATE utf8mb4_unicode_ci,
  `contractMethod` longtext COLLATE utf8mb4_unicode_ci,
  `contractParameters` longtext COLLATE utf8mb4_unicode_ci,
  `contractParametersBak` longtext COLLATE utf8mb4_unicode_ci,
  `contractVersion` longtext COLLATE utf8mb4_unicode_ci,
  `endorsement` longtext COLLATE utf8mb4_unicode_ci,
  `sequence` bigint unsigned DEFAULT NULL,
  `readSet` longtext COLLATE utf8mb4_unicode_ci,
  `readSetBak` longtext COLLATE utf8mb4_unicode_ci,
  `writeSet` longtext COLLATE utf8mb4_unicode_ci,
  `writeSetBak` longtext COLLATE utf8mb4_unicode_ci,
  `userAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `payerAddr` longtext COLLATE utf8mb4_unicode_ci,
  `gasUsed` bigint unsigned DEFAULT NULL,
  `event` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`txId`),
  KEY `timestamp_index` (`timestamp`),
  KEY `contract_name_index` (`contractNameBak`),
  KEY `contract_addr_index` (`contractAddr`),
  KEY `user_addr_index` (`userAddr`),
  KEY `sender_addr_index` (`sender`),
  KEY `block_height_index` (`blockHeight`),
  KEY `block_hash_index` (`blockHash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_block`
--

DROP TABLE IF EXISTS `chainmaker_pk_block`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_block` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `preBlockHash` longtext COLLATE utf8mb4_unicode_ci,
  `blockHash` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `blockVersion` int DEFAULT NULL,
  `orgId` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `blockDag` longtext COLLATE utf8mb4_unicode_ci,
  `dagHash` longtext COLLATE utf8mb4_unicode_ci,
  `txCount` bigint DEFAULT NULL,
  `signature` longtext COLLATE utf8mb4_unicode_ci,
  `rwSetHash` longtext COLLATE utf8mb4_unicode_ci,
  `txRootHash` longtext COLLATE utf8mb4_unicode_ci,
  `proposerId` longtext COLLATE utf8mb4_unicode_ci,
  `proposerAddr` longtext COLLATE utf8mb4_unicode_ci,
  `consensusArgs` longtext COLLATE utf8mb4_unicode_ci,
  `delayUpdateStatus` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_chainmaker_pk_block_block_height` (`blockHeight`),
  KEY `block_hash_index` (`blockHash`),
  KEY `timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_contract`
--

DROP TABLE IF EXISTS `chainmaker_pk_contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_contract` (
  `name` longtext COLLATE utf8mb4_unicode_ci,
  `nameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `addr` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `version` longtext COLLATE utf8mb4_unicode_ci,
  `runtimeType` longtext COLLATE utf8mb4_unicode_ci,
  `contractStatus` int DEFAULT NULL,
  `contractType` longtext COLLATE utf8mb4_unicode_ci,
  `contractSymbol` longtext COLLATE utf8mb4_unicode_ci,
  `decimals` bigint DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `txNum` bigint DEFAULT NULL,
  `orgId` longtext COLLATE utf8mb4_unicode_ci,
  `createTxId` longtext COLLATE utf8mb4_unicode_ci,
  `createSender` longtext COLLATE utf8mb4_unicode_ci,
  `creatorAddr` longtext COLLATE utf8mb4_unicode_ci,
  `upgradeAddr` longtext COLLATE utf8mb4_unicode_ci,
  `upgradeTimestamp` bigint DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`addr`),
  KEY `name_bak_index` (`nameBak`),
  KEY `create_timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_contract_event`
--

DROP TABLE IF EXISTS `chainmaker_pk_contract_event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_contract_event` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `eventIndex` bigint DEFAULT NULL,
  `topic` longtext COLLATE utf8mb4_unicode_ci,
  `topicBak` longtext COLLATE utf8mb4_unicode_ci,
  `contractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractNameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractVersion` longtext COLLATE utf8mb4_unicode_ci,
  `eventData` mediumblob,
  `eventDataBak` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `event_txId_index` (`txId`,`eventIndex`),
  KEY `timestamp_index` (`timestamp`),
  KEY `contract_name_index` (`contractName`),
  KEY `contract_name_bak_index` (`contractNameBak`),
  KEY `contract_addr_index` (`contractAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_contract_upgrade_transaction`
--

DROP TABLE IF EXISTS `chainmaker_pk_contract_upgrade_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_contract_upgrade_transaction` (
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `senderOrgId` longtext COLLATE utf8mb4_unicode_ci,
  `sender` longtext COLLATE utf8mb4_unicode_ci,
  `userAddr` longtext COLLATE utf8mb4_unicode_ci,
  `blockHeight` bigint unsigned DEFAULT NULL,
  `blockHash` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `txStatusCode` longtext COLLATE utf8mb4_unicode_ci,
  `contractResultCode` int unsigned DEFAULT NULL,
  `contractRuntimeType` longtext COLLATE utf8mb4_unicode_ci,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `contractNameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractVersion` longtext COLLATE utf8mb4_unicode_ci,
  `contractType` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`txId`),
  KEY `name_index` (`contractNameBak`),
  KEY `addr_index` (`contractAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_business_transaction`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_business_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_business_transaction` (
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `crossId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `subChainId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `isMainChain` tinyint(1) DEFAULT NULL,
  `gatewayId` longtext COLLATE utf8mb4_unicode_ci,
  `txStatus` int DEFAULT NULL,
  `crossContractResult` longtext COLLATE utf8mb4_unicode_ci,
  `txType` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `txStatusCode` longtext COLLATE utf8mb4_unicode_ci,
  `rwSetHash` longtext COLLATE utf8mb4_unicode_ci,
  `contractResultCode` int unsigned DEFAULT NULL,
  `contractResult` mediumblob,
  `contractMessage` blob,
  `contractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractMethod` longtext COLLATE utf8mb4_unicode_ci,
  `contractParameters` longtext COLLATE utf8mb4_unicode_ci,
  `gasUsed` bigint unsigned DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`txId`),
  KEY `rid_index` (`subChainId`),
  KEY `timestamp_index` (`timestamp`),
  KEY `contract_name_index` (`contractName`),
  KEY `cross_index` (`crossId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_chain_contract`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_chain_contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_chain_contract` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `subChainId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `sub_contract_index` (`subChainId`,`contractName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_cycle_transaction`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_cycle_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_cycle_transaction` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `crossId` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` int DEFAULT NULL,
  `startTime` bigint DEFAULT NULL,
  `endTime` bigint DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_chainmaker_pk_cross_cycle_transaction_cross_id` (`crossId`),
  KEY `start_index` (`startTime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_main_transaction`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_main_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_main_transaction` (
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `crossId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chainMsg` longtext COLLATE utf8mb4_unicode_ci,
  `status` int DEFAULT NULL,
  `crossType` int DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`txId`),
  KEY `cross_index` (`crossId`),
  KEY `timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_sub_chain_cross_chain`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_sub_chain_cross_chain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_sub_chain_cross_chain` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `subChainId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chainId` longtext COLLATE utf8mb4_unicode_ci,
  `chainName` longtext COLLATE utf8mb4_unicode_ci,
  `txNum；index:tx_num_index` bigint DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `sub_chain_index` (`subChainId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_sub_chain_data`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_sub_chain_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_sub_chain_data` (
  `subChainId` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `txNum` bigint DEFAULT NULL,
  `chainId` longtext COLLATE utf8mb4_unicode_ci,
  `chainName` longtext COLLATE utf8mb4_unicode_ci,
  `gatewayId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gatewayName` longtext COLLATE utf8mb4_unicode_ci,
  `gatewayAddr` longtext COLLATE utf8mb4_unicode_ci,
  `chainType` int DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `txVerifyType` bigint DEFAULT NULL,
  `status` int DEFAULT NULL,
  `enable` tinyint(1) DEFAULT NULL,
  `crossCa` longtext COLLATE utf8mb4_unicode_ci,
  `sdkClientCrt` longtext COLLATE utf8mb4_unicode_ci,
  `sdkClientKey` longtext COLLATE utf8mb4_unicode_ci,
  `spvContractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `introduction` longtext COLLATE utf8mb4_unicode_ci,
  `explorerAddress` longtext COLLATE utf8mb4_unicode_ci,
  `explorerTxAddress` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`subChainId`),
  KEY `gateway_index` (`gatewayId`),
  KEY `spv_name_index` (`spvContractName`),
  KEY `timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_cross_transaction_transfer`
--

DROP TABLE IF EXISTS `chainmaker_pk_cross_transaction_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_cross_transaction_transfer` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `crossId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fromGatewayId` longtext COLLATE utf8mb4_unicode_ci,
  `fromChainId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fromIsMainChain` tinyint(1) DEFAULT NULL,
  `toGatewayId` longtext COLLATE utf8mb4_unicode_ci,
  `toChainId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `toIsMainChain` tinyint(1) DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `contractMethod` longtext COLLATE utf8mb4_unicode_ci,
  `parameter` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cross_index` (`crossId`),
  KEY `from_chain_index` (`fromChainId`),
  KEY `to_chain_index` (`toChainId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_evidence_contract`
--

DROP TABLE IF EXISTS `chainmaker_pk_evidence_contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_evidence_contract` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `txId` longtext COLLATE utf8mb4_unicode_ci,
  `senderAddr` longtext COLLATE utf8mb4_unicode_ci,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `evidenceId` longtext COLLATE utf8mb4_unicode_ci,
  `hash` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `metaData` longtext COLLATE utf8mb4_unicode_ci,
  `metaDataBak` longtext COLLATE utf8mb4_unicode_ci,
  `blockHeight` bigint DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_chainmaker_pk_evidence_contract_hash` (`hash`),
  KEY `timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_fungible_contract`
--

DROP TABLE IF EXISTS `chainmaker_pk_fungible_contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_fungible_contract` (
  `symbol` longtext COLLATE utf8mb4_unicode_ci,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `contractNameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `contractType` longtext COLLATE utf8mb4_unicode_ci,
  `totalSupply` longtext COLLATE utf8mb4_unicode_ci,
  `holderCount` bigint DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`contractAddr`),
  KEY `name_bak_index` (`contractNameBak`),
  KEY `timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_fungible_position`
--

DROP TABLE IF EXISTS `chainmaker_pk_fungible_position`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_fungible_position` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `addrType` bigint DEFAULT NULL,
  `ownerAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `symbol` longtext COLLATE utf8mb4_unicode_ci,
  `amount` longtext COLLATE utf8mb4_unicode_ci,
  `blockHeight` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner_contract_index` (`ownerAddr`,`contractAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_fungible_transfer`
--

DROP TABLE IF EXISTS `chainmaker_pk_fungible_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_fungible_transfer` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `eventIndex` bigint DEFAULT NULL,
  `contractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractMethod` longtext COLLATE utf8mb4_unicode_ci,
  `topic` longtext COLLATE utf8mb4_unicode_ci,
  `fromAddr` longtext COLLATE utf8mb4_unicode_ci,
  `toAddr` longtext COLLATE utf8mb4_unicode_ci,
  `amount` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `event_txId_index` (`txId`,`eventIndex`),
  KEY `addr_index` (`contractAddr`),
  KEY `name_index` (`contractName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_gas`
--

DROP TABLE IF EXISTS `chainmaker_pk_gas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_gas` (
  `address` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `gasBalance` bigint DEFAULT NULL,
  `gasTotal` bigint DEFAULT NULL,
  `gasUsed` bigint DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_gas_record`
--

DROP TABLE IF EXISTS `chainmaker_pk_gas_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_gas_record` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `gasIndex` bigint DEFAULT NULL,
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gasAmount` bigint DEFAULT NULL,
  `businessType` bigint DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `gas_txId_index` (`gasIndex`,`txId`),
  KEY `address_index` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_identity_contract`
--

DROP TABLE IF EXISTS `chainmaker_pk_identity_contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_identity_contract` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `eventIndex` bigint DEFAULT NULL,
  `contractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userAddr` longtext COLLATE utf8mb4_unicode_ci,
  `level` longtext COLLATE utf8mb4_unicode_ci,
  `pkPem` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `event_txId_index` (`txId`,`eventIndex`),
  KEY `name_index` (`contractName`),
  KEY `addr_index` (`contractAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_node`
--

DROP TABLE IF EXISTS `chainmaker_pk_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_node` (
  `nodeId` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `nodeName` longtext COLLATE utf8mb4_unicode_ci,
  `orgId` longtext COLLATE utf8mb4_unicode_ci,
  `role` longtext COLLATE utf8mb4_unicode_ci,
  `address` longtext COLLATE utf8mb4_unicode_ci,
  `status` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`nodeId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_non_fungible_contract`
--

DROP TABLE IF EXISTS `chainmaker_pk_non_fungible_contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_non_fungible_contract` (
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `contractNameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `contractType` longtext COLLATE utf8mb4_unicode_ci,
  `totalSupply` longtext COLLATE utf8mb4_unicode_ci,
  `holderCount` bigint DEFAULT NULL,
  `blockHeight` bigint DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`contractAddr`),
  KEY `name_bak_index` (`contractNameBak`),
  KEY `timestamp_index` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_non_fungible_position`
--

DROP TABLE IF EXISTS `chainmaker_pk_non_fungible_position`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_non_fungible_position` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `addrType` bigint DEFAULT NULL,
  `ownerAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `amount` longtext COLLATE utf8mb4_unicode_ci,
  `blockHeight` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner_contract_index` (`ownerAddr`,`contractAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_non_fungible_token`
--

DROP TABLE IF EXISTS `chainmaker_pk_non_fungible_token`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_non_fungible_token` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tokenId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `addrType` bigint DEFAULT NULL,
  `ownerAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `metaData` longtext COLLATE utf8mb4_unicode_ci,
  `metaDataBak` longtext COLLATE utf8mb4_unicode_ci,
  `categoryName` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tokenId_contractAddr` (`tokenId`,`contractAddr`),
  KEY `owner_addr_index` (`ownerAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_non_fungible_transfer`
--

DROP TABLE IF EXISTS `chainmaker_pk_non_fungible_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_non_fungible_transfer` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `eventIndex` bigint DEFAULT NULL,
  `contractName` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` longtext COLLATE utf8mb4_unicode_ci,
  `contractMethod` longtext COLLATE utf8mb4_unicode_ci,
  `topic` longtext COLLATE utf8mb4_unicode_ci,
  `fromAddr` longtext COLLATE utf8mb4_unicode_ci,
  `toAddr` longtext COLLATE utf8mb4_unicode_ci,
  `tokenId` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `event_txId_index` (`txId`,`eventIndex`),
  KEY `name_index` (`contractName`),
  KEY `token_index` (`tokenId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_org`
--

DROP TABLE IF EXISTS `chainmaker_pk_org`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_org` (
  `orgId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`orgId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_transaction`
--

DROP TABLE IF EXISTS `chainmaker_pk_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_transaction` (
  `txId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sender` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `senderOrgId` longtext COLLATE utf8mb4_unicode_ci,
  `blockHeight` bigint DEFAULT NULL,
  `blockHash` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `txType` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `expirationTime` bigint DEFAULT NULL,
  `txIndex` bigint DEFAULT NULL,
  `txStatusCode` longtext COLLATE utf8mb4_unicode_ci,
  `rwSetHash` longtext COLLATE utf8mb4_unicode_ci,
  `contractResultCode` int unsigned DEFAULT NULL,
  `contractResult` mediumblob,
  `contractResultBak` mediumblob,
  `contractMessage` blob,
  `contractMessageBak` blob,
  `contractName` longtext COLLATE utf8mb4_unicode_ci,
  `contractNameBak` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contractRuntimeType` longtext COLLATE utf8mb4_unicode_ci,
  `contractType` longtext COLLATE utf8mb4_unicode_ci,
  `contractMethod` longtext COLLATE utf8mb4_unicode_ci,
  `contractParameters` longtext COLLATE utf8mb4_unicode_ci,
  `contractParametersBak` longtext COLLATE utf8mb4_unicode_ci,
  `contractVersion` longtext COLLATE utf8mb4_unicode_ci,
  `endorsement` longtext COLLATE utf8mb4_unicode_ci,
  `sequence` bigint unsigned DEFAULT NULL,
  `readSet` longtext COLLATE utf8mb4_unicode_ci,
  `readSetBak` longtext COLLATE utf8mb4_unicode_ci,
  `writeSet` longtext COLLATE utf8mb4_unicode_ci,
  `writeSetBak` longtext COLLATE utf8mb4_unicode_ci,
  `userAddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `payerAddr` longtext COLLATE utf8mb4_unicode_ci,
  `gasUsed` bigint unsigned DEFAULT NULL,
  `event` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`txId`),
  KEY `contract_addr_index` (`contractAddr`),
  KEY `user_addr_index` (`userAddr`),
  KEY `sender_addr_index` (`sender`),
  KEY `block_height_index` (`blockHeight`),
  KEY `block_hash_index` (`blockHash`),
  KEY `timestamp_index` (`timestamp`),
  KEY `contract_name_index` (`contractNameBak`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chainmaker_pk_user`
--

DROP TABLE IF EXISTS `chainmaker_pk_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chainmaker_pk_user` (
  `userId` longtext COLLATE utf8mb4_unicode_ci,
  `userAddr` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` longtext COLLATE utf8mb4_unicode_ci,
  `orgId` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`userAddr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chain`
--

DROP TABLE IF EXISTS `chain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chain` (
  `chainId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `version` longtext COLLATE utf8mb4_unicode_ci,
  `chainName` longtext COLLATE utf8mb4_unicode_ci,
  `enableGas` tinyint(1) DEFAULT NULL,
  `consensus` longtext COLLATE utf8mb4_unicode_ci,
  `txTimestampVerify` tinyint(1) DEFAULT NULL,
  `txTimeout` bigint DEFAULT NULL,
  `blockTxCapacity` bigint DEFAULT NULL,
  `blockSize` bigint DEFAULT NULL,
  `blockInterval` bigint DEFAULT NULL,
  `hashType` longtext COLLATE utf8mb4_unicode_ci,
  `authType` longtext COLLATE utf8mb4_unicode_ci,
  `timestamp` bigint DEFAULT NULL,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`chainId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `subscribe`
--

DROP TABLE IF EXISTS `subscribe`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `subscribe` (
  `chainId` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `orgId` longtext COLLATE utf8mb4_unicode_ci,
  `userKey` text COLLATE utf8mb4_unicode_ci,
  `userCert` text COLLATE utf8mb4_unicode_ci,
  `nodeList` text COLLATE utf8mb4_unicode_ci,
  `status` bigint DEFAULT NULL,
  `authType` longtext COLLATE utf8mb4_unicode_ci,
  `hashType` longtext COLLATE utf8mb4_unicode_ci,
  `nodeCACert` text COLLATE utf8mb4_unicode_ci,
  `tls` tinyint(1) DEFAULT NULL,
  `tlsHost` longtext COLLATE utf8mb4_unicode_ci,
  `remote` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) DEFAULT NULL,
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`chainId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-03-22 15:46:11
