阶段：同步数据
====
从tb表中读取数据，按照路由规则，写入tb_n表中：
1.自增ID主键倒排分页批量取数据
2.如果数据在tb_n中不存在，则插入
3.如果数据在tb_n中存在，且tb中set_time大于tb_n中set_time，则覆盖

INSERT INTO tb_0 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 0;
INSERT INTO tb_1 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 1;
INSERT INTO tb_2 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 2;
INSERT INTO tb_3 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 3;
INSERT INTO tb_4 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 4;
INSERT INTO tb_5 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 5;
INSERT INTO tb_6 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 6;
INSERT INTO tb_7 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 7;
INSERT INTO tb_8 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 8;
INSERT INTO tb_9 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 9;
INSERT INTO tb_10 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 10;
INSERT INTO tb_11 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 11;
INSERT INTO tb_12 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 12;
INSERT INTO tb_13 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 13;
INSERT INTO tb_14 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 14;
INSERT INTO tb_15 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 15;
INSERT INTO tb_16 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 16;
INSERT INTO tb_17 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 17;
INSERT INTO tb_18 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 18;
INSERT INTO tb_19 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 19;
INSERT INTO tb_20 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 20;
INSERT INTO tb_21 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 21;
INSERT INTO tb_22 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 22;
INSERT INTO tb_23 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 23;
INSERT INTO tb_24 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 24;
INSERT INTO tb_25 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 25;
INSERT INTO tb_26 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 26;
INSERT INTO tb_27 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 27;
INSERT INTO tb_28 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 28;
INSERT INTO tb_29 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 29;
INSERT INTO tb_30 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 30;
INSERT INTO tb_31 SELECT * FROM tb WHERE kid % (1 * 32) % 32 = 31;






