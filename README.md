# Project - Backend API

The goal of this project is to create an http+json API for a calculator service.

## Overview

This calculator is stateless, meaning that there is no data stored in a database or in memory.

## Features

- Input Validation
- Error feedback
- Logging

## API Endpoints

### Add Two Numbers

- **URL**: `/add`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "number1": 5,
    "number2": 3
  }
  ```
- **Response**:
  ```json
  {
    "result": 8
  }
  ```

### Subtract Two Numbers

- **URL**: `/subtract`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "number1": 5,
    "number2": 3
  }
  ```
- **Response**:
  ```json
  {
    "result": 2
  }
  ```

### Multiply Two Numbers

- **URL**: `/multiply`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "number1": 5,
    "number2": 3
  }
  ```
- **Response**:
  ```json
  {
    "result": 15
  }
  ```

### Divide Two Numbers

- **URL**: `/divide`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "number1": 6,
    "number2": 3
  }
  ```
- **Response**:
  ```json
  {
    "result": 2
  }
  ```

### Add All Numbers in an Array

- **URL**: `/sum`
- **Method**: `POST`
- **Request Body**:
  ```json
  [1, 2, 3, 4, 5]
  ```
- **Response**:
  ```json
  {
    "result": 15
  }
  ```

## Additional Tasks

- [x] Add in rate limiter to prevent misuse of the API
- [ ] Add in token authentication to prevent anyone unauthorized from using the API
- [ ] Add in a database to keep track of all of the calculations that have taken place
- [ ] Add in support for floating point numbers as well.
- [ ] Create an associated http client that can work with the calculator API.
- [ ] Create a frontend that makes use of your API.
- [ ] Add in a middleware that adds a request ID to the http.Request object.
