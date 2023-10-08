/* metadata table for product data */
CREATE TABLE product_meta (
  id          varchar(32) NOT NULL, /* an identifier form : {type}.{name}.{location} */
  name        varchar(32) NOT NULL, /* human readable stock name */
  symbol      varchar(8)  NOT NULL, /* stock symbol */
  description text,                 /* product description : can be gain by external apis, or written manually by admin */
  type        varchar(8)  NOT NULL, /* examples are stock, encrypt */
  exchange    varchar(32) NOT NULL, /* examples are kospi, nasdaq. */
  location    varchar(32),          /* examples are kor, usa. when product type is coin location is null*/
  PRIMARY KEY (id)
);

CREATE TABLE platform (
  name        varchar(32) NOT NULL, /* available platform is buycycle, polygon, kis */
  description text,
  PRIMARY KEY (name)
);

CREATE TABLE product_platform (
  product_id    varchar(32) NOT NULL,
  platform_name varchar(32) NOT NULL, 
  identifier    varchar(32) NOT NULL, /* a string that is used to specific stock on such platform query */

  PRIMARY KEY (product_id, platform_name),
  FOREIGN KEY (product_id) REFERENCES product_meta (id),
  FOREIGN KEY (platform_name) REFERENCES platform (name)
);

CREATE TABLE store_log (
  product_id  varchar(32) NOT NULL,
  stored_at   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "status"    varchar(10) NOT NULL,

  PRIMARY KEY (stored_at),
  FOREIGN KEY (product_id) REFERENCES product_meta (id)
);
