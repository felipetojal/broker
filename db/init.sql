
CREATE TABLE  users
(                           --USUARIOS
    id UUID DEFAULT gen_random_uuid(),
    username   UNIQUE NOT NULL,
    senha_hash TEXT  IS NOT NULL, --as senhas nao precisam ser diferentes
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE accounts
(                                                                  --CONTAS
    UUID PRIMARY KEY REFERENCES users(id),                             --nao pode ter serial pq migraremos para outro DB mais na frente
    saldo BIGINT   NOT NULL CHECK(saldo>=0),
    reservado BIGINT  NOT NULL DEFAULT 0 CHECK (reservado >= 0 AND reservado <= saldo),--como pego BIGINT? Aqui é em centavos
    FOREIGN KEY (id) REFERENCES users (id) SET NULL ON UPDATE CASCADE
);

--ACOES
CREATE TABLE assets
(
    --preciso adicionar um ticker que se chama assets, no caso a acao que vai ter silga(chave), nome, setor, preco atual em centavos e por ai vai
    ticker TEXT PRIMARY KEY  NOT NULL, --é a sigla
    nome TEXT UNIQUE  NOT NULL,
    setor TEXT  NOT NULL,
    valor BIGINT CHECK (valor >= 0.0) --acao com centavos, verificar depois
);


--NOME DAS EMPRESAS DOS DOCUMENTOS
INSERT INTO assets VALUES('ARA','aracaju','publico',20.9);
INSERT INTO assets VALUES('XLA','Xlag','espacial',45.3);
INSERT INTO assets VALUES('SOS','Sosoft','saude', 10.2);
INSERT INTO assets VALUES('APR','Appria', 'tecnoliga',145.2);
INSERT INTO assets VALUES('VOK','Volkjeira','automobilistica',30.3);


--orders
CREATE TABLE orders
(
    id_compra UUID DEFAULT gen_random_uuid(),    --É PRIMARY KEY? DEVVE SER COMO SERIAL, VEJA!
    UUID PRIMARY KEY REFERENCES users(id),
    ticker TEXT NOT NULL REFERENCES assets(ticker),
    lado CHAR[5]    CHECK (lado IN ('BUY', 'SELL')), --Ou compra ou vende
    preco BIGINT NOT NULL CHECK (preco>0),                    --0
    quantidade INT CHECK (quantidade>0),
    status TEXT CHECK(status IN('OPEN','PARTIAL','FILLED','CANCELLED','EXPIRED','REJECTED')),
    criada_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    quantidade_restante INT,
    CREATE INDEX indice ON orders (ticker,status)
    FOREIGN KEY (id) REFERENCES users(id)  SET NULL ON UPDATE CASCADE

    --veja se deve ser direto do accounts ou do users
);



--positions
CREATE TABLE positions
(
    PRIMARY KEY(id,ticker) REFERENCES user(id),
    quantidade INT    , --TALVEZ SEJA EM BIGINT, para cripto por exemplo
    reservada INT NOT NULL DEFAULT 0 CHECK(reservada>=0 AND reservada<=quantidade ),         --reservada vai ser o valor ou um estado(sim ou nao)?
    preco_medio BIGINT
);


--trades
CREATE TABLE trades
(
    --deve ter o ID da trade unico,user id,valor da transacao,atualizar o saldo da conta FG com accounts
    --, nome e sigla da acao FG de assets, quantidade, feita em
    id_trade INT     PRIMARY KEY  NOT NULL, --mesma ideia do SERIAL, ver como vai fazer
    user_id DEFAULT gen_random_uuid(),       --mesma ideia, vai pegar de onde?
    valor_transacao BIGINT ,
    atualizar_saldo BIGINT   NOT NULL,      --FG, pego o novo valor apos a transcacao e atualizo aconta
    nome_acao TEXT     NOT NULL,            --FG
    sigla TEXT     NOT NULL,--FG
    quantidade INT    ,
    data_transacao TIMESTAMPTZ NOT NULL DEFAULT now()   ,
    FOREIGN KEY (user_id) REFERENCES accounts (user_id) ON DELETE SET NULL ON UPDATE CASCATE,
    FOREIGN KEY (atualizar_saldo) REFERENCES accounts (saldo) ON DELETE SET NULL ON UPDATE CASCATE,
    FOREIGN KEY (nome_acao) REFERENCES assets (nome_acao) ON DELETE SET NULL ON UPDATE CASCATE,
    FOREIGN KEY (sigla) REFERENCES assets (sigla) ON DELETE SET NULL ON UPDATE CASCATE

);