-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Dec 18, 2024 at 02:41 PM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `waterfalls`
--

-- --------------------------------------------------------

--
-- Table structure for table `accounts`
--

CREATE TABLE `accounts` (
  `Id` int(11) NOT NULL,
  `FirstName` varchar(255) NOT NULL,
  `LastName` varchar(255) NOT NULL,
  `Email` varchar(255) NOT NULL,
  `Area` varchar(255) NOT NULL,
  `Password` varchar(255) NOT NULL,
  `QRCode` text NOT NULL,
  `Role` enum('Admin','Staff','Customer') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `accounts`
--

INSERT INTO `accounts` (`Id`, `FirstName`, `LastName`, `Email`, `Area`, `Password`, `QRCode`, `Role`) VALUES
(1, 'Zeke', 'Zeke', 'Zeke', 'Zeke', 'Zeke', 'Zeke', 'Customer'),
(6, 'Edison', 'Pagatpat', 'pagatpatedison@gmail.com', 'Cayang, Bogo City, Cebu', 'edison123', '1234567', 'Staff'),
(28, 'Ezekiel Angelo', 'Pelayo', 'pelayoezekiel96@gmail.com', '', '123', '', 'Admin'),
(29, 'Levi Jay', 'Pelayo', 'levi@mail.com', 'Guadalupe, Bogo City, Cebu', '123', 'levi', 'Customer'),
(30, 'qweqw', 'eweqw', 'qeqwe@mail.com', 'Guadalupe, Bogo City, Cebu', '123', 'qwe', 'Customer'),
(31, 'qweqwe', 'qweqwe', 'qwew@mail.com', 'qwewq', '123', 'weq', 'Customer'),
(32, 'Ezekiel Angelo', 'Pelayo', 'eianezekiel@yahoo.com', 'Cayang, Bogo City, Cebu', '123', 'qwe', 'Customer'),
(33, 'Pengwin', 'Kobayashi', 'pengwinkobayashi@gmail.com', 'Gairan, Bogo City, Cebu', '1234567890', 'jsdalfhasjdfhjaskfd2414', 'Admin');

-- --------------------------------------------------------

--
-- Table structure for table `agents`
--

CREATE TABLE `agents` (
  `Id` int(11) NOT NULL,
  `area_id` int(11) NOT NULL,
  `agent_name` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `agents`
--

INSERT INTO `agents` (`Id`, `area_id`, `agent_name`) VALUES
(90, 10, 'qweqwe'),
(91, 12, 'w'),
(92, 12, 'eqw'),
(93, 12, 'qweqweqwewq'),
(94, 10, 'w'),
(95, 10, 'qwe'),
(96, 0, 'weqe'),
(97, 0, 'qwe'),
(98, 10, '2');

-- --------------------------------------------------------

--
-- Table structure for table `areas`
--

CREATE TABLE `areas` (
  `Id` int(11) NOT NULL,
  `Area` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `areas`
--

INSERT INTO `areas` (`Id`, `Area`) VALUES
(10, 'Guadalupe, Bogo City, Cebu'),
(12, 'Nailon, Bogo City, Cebu'),
(13, 'Lapaz, Bogo City, Cebu'),
(14, 'Malingin, Bogo City, Cebu'),
(18, 'qwe');

-- --------------------------------------------------------

--
-- Table structure for table `customer_order`
--

CREATE TABLE `customer_order` (
  `Id` int(11) NOT NULL,
  `num_gallons_order` int(11) NOT NULL,
  `date` varchar(255) NOT NULL,
  `date_created` timestamp NOT NULL DEFAULT current_timestamp(),
  `customer_id` int(11) DEFAULT NULL,
  `total_price` decimal(11,2) NOT NULL,
  `status` varchar(200) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `customer_order`
--

INSERT INTO `customer_order` (`Id`, `num_gallons_order`, `date`, `date_created`, `customer_id`, `total_price`, `status`) VALUES
(2, 10, '20241212', '2024-12-12 06:35:11', NULL, 0.00, ''),
(3, 10, 'Thursday', '2024-12-12 06:36:36', NULL, 0.00, ''),
(4, 100, 'Tuesday', '2024-12-12 08:15:15', NULL, 0.00, ''),
(5, 10, 'Tuesday', '2024-12-12 08:18:17', NULL, 0.00, ''),
(6, 10, 'Tuesday', '2024-12-12 08:18:19', NULL, 0.00, ''),
(7, 10, 'Tuesday', '2024-12-12 08:20:13', NULL, 0.00, ''),
(8, 10, 'Tuesday', '2024-12-12 08:20:57', NULL, 0.00, ''),
(9, 10, 'Tuesday', '2024-12-12 08:21:08', NULL, 0.00, ''),
(10, 10, 'Tuesday', '2024-12-12 08:21:11', NULL, 0.00, ''),
(11, 10, 'Tuesday', '2024-12-12 08:21:40', NULL, 0.00, ''),
(12, 20, 'Tuesday', '2024-12-12 08:22:20', NULL, 0.00, ''),
(13, 20, 'Tuesday', '2024-12-12 08:25:49', NULL, 0.00, ''),
(14, 1, 'Thursday', '2024-12-12 08:41:06', NULL, 0.00, ''),
(15, 2, 'Thursday', '2024-12-12 08:50:52', NULL, 0.00, ''),
(16, 2, 'Thursday', '2024-12-12 08:50:57', NULL, 0.00, ''),
(17, 14, 'Thursday', '2024-12-12 08:51:15', NULL, 0.00, ''),
(18, 20, 'Tuesday', '2024-12-12 08:53:28', NULL, 0.00, ''),
(19, 12, 'Tuesday', '2024-12-13 07:15:32', NULL, 0.00, ''),
(20, 100, 'Thursday', '2024-12-13 07:57:24', NULL, 0.00, ''),
(21, 1, 'Thursday', '2024-12-13 09:15:40', NULL, 0.00, ''),
(22, 21, 'Tuesday', '2024-12-13 09:20:49', NULL, 0.00, ''),
(23, 1, 'Thursday', '2024-12-17 09:01:42', NULL, 0.00, ''),
(24, 100, 'Thursday', '2024-12-18 06:30:01', NULL, 0.00, ''),
(25, 69, 'Tuesday', '2024-12-18 06:30:41', NULL, 0.00, ''),
(26, 111, 'Tuesday', '2024-12-18 06:50:22', NULL, 0.00, ''),
(27, 10, '2024-12-18', '2024-12-18 07:06:21', 32, 0.00, ''),
(28, 555, '2024-12-18', '2024-12-18 09:14:30', 32, 0.00, ''),
(29, 123, 'Tuesday', '2024-12-18 09:41:30', 32, 0.00, ''),
(30, 666999, 'Tuesday', '2024-12-18 09:43:14', 32, 0.00, '');

-- --------------------------------------------------------

--
-- Table structure for table `inventory`
--

CREATE TABLE `inventory` (
  `inventory_id` int(11) NOT NULL,
  `item` varchar(255) NOT NULL,
  `no_of_items` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `inventory`
--

INSERT INTO `inventory` (`inventory_id`, `item`, `no_of_items`) VALUES
(1, 'q', 100),
(49, 'qewqweqw', 32);

-- --------------------------------------------------------

--
-- Table structure for table `inventory_available`
--

CREATE TABLE `inventory_available` (
  `inventory_id` int(11) NOT NULL,
  `total_quantity` int(11) NOT NULL,
  `price` decimal(11,2) NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `messages`
--

CREATE TABLE `messages` (
  `id` int(11) NOT NULL,
  `sender` varchar(255) DEFAULT NULL,
  `recipient` varchar(255) DEFAULT NULL,
  `content` text DEFAULT NULL,
  `timestamp` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `messages`
--

INSERT INTO `messages` (`id`, `sender`, `recipient`, `content`, `timestamp`) VALUES
(1, 'user', 'Agent', 'hi', '2024-12-04 02:42:28'),
(2, 'agent', 'Ricky Monsales', 'qwe', '2024-12-04 02:42:50'),
(3, 'user', 'Agent', 'hi', '2024-12-04 02:45:45'),
(4, 'agent', 'Ricky Monsales', 'yow', '2024-12-04 02:47:33'),
(5, 'user', 'Agent', 'ewq', '2024-12-04 02:48:18'),
(6, 'agent', 'Admin', 'wewq', '2024-12-04 02:49:16'),
(7, 'user', 'Agent', 'hi\\', '2024-12-04 02:53:54'),
(8, 'agent', 'Admin', 'yes', '2024-12-04 02:54:08'),
(9, 'agent', 'Ricky Monsales', 'hshs', '2024-12-04 06:43:01'),
(10, 'agent', 'Admin', 'sir', '2024-12-04 06:43:08'),
(11, 'user', 'Agent', 'sir', '2024-12-04 06:48:50'),
(12, 'user', 'Agent', 'gdfgdfgd', '2024-12-06 12:31:53'),
(13, 'user', 'Agent', 'rttrtr', '2024-12-06 12:32:33'),
(14, 'user', 'Agent', 'ffdsf', '2024-12-06 12:42:31'),
(15, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:51:29'),
(16, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:51:49'),
(17, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:51:51'),
(18, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:01'),
(19, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:05'),
(20, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:06'),
(21, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:13'),
(22, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:14'),
(23, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:15'),
(24, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:16'),
(25, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:27'),
(26, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:27'),
(27, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:28'),
(28, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:28'),
(29, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:29'),
(30, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:29'),
(31, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:30'),
(32, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:30'),
(33, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:31'),
(34, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:31'),
(35, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:32'),
(36, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:32'),
(37, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:32'),
(38, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:32'),
(39, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:33'),
(40, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:33'),
(41, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:33'),
(42, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:33'),
(43, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:33'),
(44, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:33'),
(45, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:34'),
(46, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:34'),
(47, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:34'),
(48, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:34'),
(49, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:34'),
(50, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:37'),
(51, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:38'),
(52, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:38'),
(53, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:38'),
(54, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:38'),
(55, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:39'),
(56, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:39'),
(57, 'Hello', 'Hello', 'Heeloo', '2024-12-06 12:52:39');

-- --------------------------------------------------------

--
-- Table structure for table `staffs`
--

CREATE TABLE `staffs` (
  `id` int(11) NOT NULL,
  `staff_name` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `staffs`
--

INSERT INTO `staffs` (`id`, `staff_name`, `address`) VALUES
(16, 'qweqwe', 'qweqwe');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `accounts`
--
ALTER TABLE `accounts`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `agents`
--
ALTER TABLE `agents`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `areas`
--
ALTER TABLE `areas`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `customer_order`
--
ALTER TABLE `customer_order`
  ADD PRIMARY KEY (`Id`);

--
-- Indexes for table `inventory`
--
ALTER TABLE `inventory`
  ADD PRIMARY KEY (`inventory_id`);

--
-- Indexes for table `inventory_available`
--
ALTER TABLE `inventory_available`
  ADD PRIMARY KEY (`inventory_id`);

--
-- Indexes for table `messages`
--
ALTER TABLE `messages`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `staffs`
--
ALTER TABLE `staffs`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `accounts`
--
ALTER TABLE `accounts`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=34;

--
-- AUTO_INCREMENT for table `agents`
--
ALTER TABLE `agents`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=99;

--
-- AUTO_INCREMENT for table `areas`
--
ALTER TABLE `areas`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=19;

--
-- AUTO_INCREMENT for table `customer_order`
--
ALTER TABLE `customer_order`
  MODIFY `Id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=31;

--
-- AUTO_INCREMENT for table `inventory`
--
ALTER TABLE `inventory`
  MODIFY `inventory_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=50;

--
-- AUTO_INCREMENT for table `inventory_available`
--
ALTER TABLE `inventory_available`
  MODIFY `inventory_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `messages`
--
ALTER TABLE `messages`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=58;

--
-- AUTO_INCREMENT for table `staffs`
--
ALTER TABLE `staffs`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
