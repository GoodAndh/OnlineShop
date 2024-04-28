create table if not exists orders_items(
    id int not null auto_increment,
    total int not null,
    total_price float not null,
    order_id int not null ,
    product_id int not null ,
    users_id int not null ,
    foreign key (users_id)references users(id),
    foreign key (order_id) references orders(id),
    foreign key (product_id) references product(id),
    primary key(id)
)engine=innodb