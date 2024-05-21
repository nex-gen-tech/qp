# Guide for Frontend Team: Writing and Passing URL Queries with QP Package

## Introduction

This document aims to help the frontend team effectively write and pass URL queries using the QP package. By following this guide, you will learn to construct dynamic queries that allow for filtering, sorting, field selection, and more. We will demonstrate these concepts with a sample table and provide 30 medium to complex queries for various use cases.

## Sample Table

Let's assume we have a table called `products` with the following schema:

```sql
CREATE TABLE products (
    product_id INT PRIMARY KEY,
    product_name VARCHAR(255),
    category VARCHAR(100),
    price DECIMAL(10, 2),
    in_stock BOOLEAN,
    created_at DATE,
    updated_at DATE
);
```

## Example Queries

Here are 30 medium to complex queries demonstrating different use cases:

### 1. Select all products with default settings
```
URL: /products
```

### 2. Select specific fields: `product_id` and `product_name`
```
URL: /products?fields=product_id,product_name
```

### 3. Sort products by `price` in ascending order
```
URL: /products?sort=price
```

### 4. Sort products by `price` descending and `product_name` ascending
```
URL: /products?sort=-price,product_name
```

### 5. Limit the results to 10 products
```
URL: /products?limit=10
```

### 6. Offset the results by 5 products
```
URL: /products?offset=5
```

### 7. Filter products by `category` equal to "electronics"
```
URL: /products?category[eq]=electronics
```

### 8. Filter products with `price` greater than 100
```
URL: /products?price[gt]=100
```

### 9. Filter products with `price` between 50 and 150
```
URL: /products?price[gte]=50&price[lte]=150
```

### 10. Filter products in stock
```
URL: /products?in_stock[eq]=true
```

### 11. Filter products out of stock
```
URL: /products?in_stock[eq]=false
```

### 12. Search for products with `product_name` containing "smart"
```
URL: /products?product_name[like]=*smart*
```

### 13. Filter products by multiple categories
```
URL: /products?category[in]=electronics,appliances
```

### 14. Exclude products from specific categories
```
URL: /products?category[nin]=furniture,clothing
```

### 15. Select products created on a specific date
```
URL: /products?created_at[eq]=2024-01-01
```

### 16. Select products updated after a specific date
```
URL: /products?updated_at[gt]=2024-01-01
```

### 17. Complex filter: products with `price` less than 100 or in the "sale" category
```
URL: /products?price[lt]=100|category[eq]=sale
```

### 18. Multiple OR conditions: name contains "smart" or category is "electronics"
```
URL: /products?product_name[like]=*smart*|category[eq]=electronics
```

### 19. Products not updated in the last 30 days
```
URL: /products?updated_at[lt]=2024-04-21
```

### 20. Products created in the last week
```
URL: /products?created_at[gte]=2024-05-14
```

### 21. Products with specific `product_id`s
```
URL: /products?product_id[in]=1,2,3,4,5
```

### 22. Products sorted by multiple fields with limit and offset
```
URL: /products?sort=price,-created_at&limit=10&offset=20
```

### 23. Products with both `product_name` and `category` matching specific patterns
```
URL: /products?product_name[like]=*smart*&category[like]=*tech*
```

### 24. Products with price not equal to 99.99
```
URL: /products?price[ne]=99.99
```

### 25. Products where `product_name` is NULL
```
URL: /products?product_name[is]=null
```

### 26. Products where `product_name` is NOT NULL
```
URL: /products?product_name[not]=null
```

### 27. Combined complex query: fields, sorting, limit, offset, and filters
```
URL: /products?fields=product_id,product_name,price&sort=-price&limit=5&offset=10&price[gt]=50&in_stock[eq]=true
```

### 28. Filter products with name starting with "Pro"
```
URL: /products?product_name[like]=Pro*
```

### 29. Select fields and filter by `created_at` range
```
URL: /products?fields=product_id,product_name,created_at&created_at[gte]=2024-01-01&created_at[lte]=2024-12-31
```

### 30. Dynamic filtering by query string from frontend inputs
```
URL: /products?fields=product_id,product_name,price,category&sort=price&limit=20&offset=0&category[in]=electronics,appliances&price[gt]=100&in_stock[eq]=true
```

## Conclusion

By using the QP package, you can construct complex and dynamic SQL queries easily through URL parameters. This guide provides examples to get you started with various filtering, sorting, and field selection scenarios. Experiment with these examples to fit your specific use cases and enhance the querying capabilities of your web applications.