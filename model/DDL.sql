
CREATE TABLE Users (
    User_ID SERIAL PRIMARY KEY,
    Name VARCHAR(50) NOT NULL,
    Email VARCHAR(50) NOT NULL,
    Password VARCHAR(50) NOT NULL,
    Phone_Number VARCHAR(20) NOT NULL,
    Address VARCHAR(100) NOT NULL,
    Balance DECIMAL(10, 2) NOT NULL
);

CREATE TABLE Transaction (
    Transaction_ID SERIAL PRIMARY KEY,
    User_ID INT NOT NULL,
    Type VARCHAR(10) NOT NULL,
    Amount DECIMAL(10, 2) NOT NULL,
    Date DATE NOT NULL,
    Description VARCHAR(100),
    FOREIGN KEY (User_ID) REFERENCES Users(User_ID)
);

CREATE TABLE Bank_Account (
    Account_ID SERIAL PRIMARY KEY,
    User_ID INT NOT NULL,
    Bank_Name VARCHAR(50) NOT NULL,
    Account_Number VARCHAR(20) NOT NULL,
    Account_Holder_Name VARCHAR(50) NOT NULL,
    FOREIGN KEY (User_ID) REFERENCES Users(User_ID)
);

CREATE TABLE Card (
    Card_ID SERIAL PRIMARY KEY,
    User_ID INT NOT NULL,
    Card_Type VARCHAR(10) NOT NULL,
    Card_Number VARCHAR(20) NOT NULL,
    Expiration_Date DATE NOT NULL,
    CVV VARCHAR(5) NOT NULL,
    FOREIGN KEY (User_ID) REFERENCES Users(User_ID)
);