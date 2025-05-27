import { JSX } from 'solid-js';
import { useNavigate } from '@solidjs/router'

import './TermoDeAceite.css';

const TermoDeAceite = ():JSX.Element => {
  const navigate = useNavigate();
  
  const handleProsseguir = () => {
    navigate('/sendMsg')
  }

  return (
    <div class="termo-container">
      <h1>Canal de Confiança</h1>

      <p>
        A Similar Tecnologia e Automação LTDA mantém o compromisso de preservar um ambiente de
        trabalho íntegro e práticas empresariais alinhadas aos seus princípios éticos, valores
        organizacionais e à legislação vigente. Por esse motivo, disponibilizamos a todos os
        colaboradores em geral este <strong>Canal de Confiança</strong> — um meio de comunicação
        independente, seguro e imparcial, destinado ao relato de denúncias anônimas sobre condutas
        que violem ou contrariem as diretrizes estabelecidas no nosso {' '}
        <a href="/manual-de-procedimentos.docx" target="_blank" rel="noopener noreferrer">
            manual de procedimentos
        </a>.
      </p>

      <h2>Importante:</h2>

      <p>
        Ao registrar uma manifestação, você declara estar ciente de que os dados informados serão
        tratados com confidencialidade, exclusivamente para fins de apuração e resolução da
        manifestação, nos termos da LGPD (Lei Geral de Proteção de Dados Pessoais – Lei nº
        13.709/2018).
      </p>

      <p>
        É fundamental que as informações fornecidas sejam verdadeiras e precisas, para que a
        apuração seja realizada de forma adequada e justa.
      </p>

      <p>
        Caso, após a devida apuração, fique constatado que a manifestação registrada é
        intencionalmente falsa ou feita de má-fé, medidas cabíveis serão tomadas.
      </p>

      <h2>Como atuamos:</h2>

      <p>
        Todas as manifestações são recebidas e analisadas com independência, isenção e absoluto
        sigilo pelo Comitê de Ética.
      </p>

      <p>Entre as principais situações que podem ser reportadas, destacam-se:</p>

      <ul>
        <li>Conflito de interesses;</li>
        <li>Assédio moral ou sexual;</li>
        <li>Irregularidades e delitos;</li>
        <li>Desrespeito às normas de segurança;</li>
        <li>Condutas antiéticas no ambiente de trabalho;</li>
        <li>Atos discriminatórios, independentemente de sua natureza;</li>
      </ul>

      <p>
        Sempre que possível e cabível, buscamos identificar alternativas para solucionar as
        demandas, com base nas informações disponíveis.
      </p>

      <p>
        O status da apuração pode ser acompanhado pelo denunciante, por meio deste canal, de forma
        objetiva e rápida.
      </p>

      <p>
        Agradecemos pelo uso responsável deste canal e pela sua colaboração com a promoção de um
        ambiente ético, íntegro e seguro na Similar.
      </p>

      <p class="termo-conclusao">
        Ao prosseguir, você declara ter lido, compreendido e aceito integralmente este Termo,
        estando ciente sobre o tratamento dos dados informados e a importância da veracidade das
        informações prestadas.
      </p>

      <button class="botao-prosseguir" onClick={handleProsseguir}>
        Prosseguir
      </button>
    </div>
  );
};

export default TermoDeAceite;