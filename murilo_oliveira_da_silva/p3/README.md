# Calculadora Brainfuck

Este projeto contém dois programas escritos em Go:

- **bfc**: Compilador que converte expressões aritméticas simples (ex: `CRÉDITO=10+2`) em código Brainfuck.
- **bfe**: Interpretador de Brainfuck que executa o código gerado.

## Como compilar

No diretório `murilo_oliveira_da_silva/p3`, execute:

```sh
make
```

Isso irá gerar executáveis bfc e bfe

## Como usar:

Rodar no terminal:

```sh
echo 'CRÉDITO=SUA-EXPRESSÃO' | ./bfc | ./bfe
```

## TODO
- Adicionar a divisão.