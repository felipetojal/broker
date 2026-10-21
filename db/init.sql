CREATE TABLE IF NOT EXISTS;

CREATE TABLE users { --USUARIOS
    UUID DEFAULT gen_random_uuid() ,
    CHAR username PRIMARY KEY IS NOT NULL,
    CHAR senha IS NOT NULL, --as senhas nao precisam ser diferentes
    DATE criado_em
};

CREATE TABLE accounts{ --CONTAS
    user_id DEFAULT gen_random_uuid(), --nao pode ter serial pq migraremos para outro DB mais na frente
    BIGINT saldo IS NOT NULL,
    BIGINT reservado CHECK (reservado>=0 AND reservado <=saldo),--como pego BIGINT? Aqui é em centavos
    FOREIGN KEY (user_id) REFERENCES users(UUID) ON DELETE SET NULL ON UPDATE CASCATE
};

--ACOES
CREATE TABLE assets{ #--preciso adicionar um ticker que se chama assets, no caso a acao que vai ter silga(chave), nome, setor, preco atual em centavos e por ai vai
    CHAR sigla PRIMARY KEY IS NOT NULL, --VEJA OQ É TICKER
    CHAR nome_acao PRIMARY KEY IS NOT NULL,
    CHAR setor IS NOT NULL,
    BIGINT valor CHECK (valor >= 0.0) --acao com centavos, verificar depois
};


--NOME DAS EMPRESAS DOS DOCUMENTOS
INSERT INTO assets VALUES('ARA','aracaju','publico',20.9);
INSERT INTO assets VALUES('XLA','Xlag','espacial',45.3);
INSERT INTO assets VALUES('SOS','Sosoft','saude', 10.2);
INSERT INTO assets VALUES('APR','Appria', 'tecnoliga',145.2);
INSERT INTO assets VALUES('VOK','Volkjeira','automobilistica',30.3);


--orders
CREATE TABLE orders{
    INT id_compra PRIMARY KEY IS NOT NULL, --É PRIMARY KEY? DEVVE SER COMO SERIAL, VEJA!
    user_id DEFAULT gen_random_uuid(),
    --adicionar um ticker, PROCURE!
    CHAR lado CHECK(lado IN ('BUY','SELL')), --Ou compra ou vende
    BIGINT preco IS NOT NULL, --0
    INT quantidade,
    CHAR status,
    DATE criada_em,
    FOREIGN KEY (user_id) REFERENCES accounts(user_id )ON DELETE SET NULL ON UPDATE CASCATE

        --veja se deve ser direto do accounts ou do users
};



--positions
CREATE TABLE positions{
    --user id + ticker como chave,
    INT quantidade, --TALVEZ SEJA EM BIGINT, para cripto por exemplo
    reservada, --reservada vai ser o valor ou um estado(sim ou nao)?
    BIGINT preco_medio
};


--trades
CREATE TABLE trades{
    --deve ter o ID da trade unico,user id,valor da transacao,atualizar o saldo da conta FG com accounts
    --, nome e sigla da acao FG de assets, quantidade, feita em
    INT id_trade PRIMARY KEY IS NOT NULL, --mesma ideia do SERIAL, ver como vai fazer
    user_id DEFAULT gen_random_uuid(), --mesma ideia, vai pegar de onde?
    BIGINT valor_transacao,
    BIGINT atualizar_saldo IS NOT NULL, --FG, pego o novo valor apos a transcacao e atualizo aconta
    CHAR nome_acao IS NOT NULL, --FG
    CHAR sigla IS NOT NULL,--FG
    INT quantidade,
    DATE data_transacao,
    FOREIGN KEY (user_id) REFERENCES accounts(user_id )ON DELETE SET NULL ON UPDATE CASCATE,
    FOREIGN KEY (atualizar_saldo) REFERENCES accounts(saldo) ON DELETE SET NULL ON UPDATE CASCATE,
    FOREIGN KEY (nome_acao) REFERENCES assets(nome_acao) ON DELETE SET NULL ON UPDATE CASCATE,
    FOREIGN KEY (sigla) REFERENCES assets(sigla) ON DELETE SET NULL ON UPDATE CASCATE

};