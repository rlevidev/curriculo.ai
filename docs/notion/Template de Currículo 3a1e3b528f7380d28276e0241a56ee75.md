# Template de Currículo

```latex
%!TEX program = pdflatex
\documentclass[11pt,a4paper]{article}

% Packages
\usepackage[T1]{fontenc}
\usepackage[left=0.6in,right=0.6in,top=0.45in,bottom=0.45in]{geometry}
\usepackage{array}
\usepackage{xcolor}
\usepackage{hyperref}
\usepackage{enumitem}
\usepackage{titlesec}
\usepackage[utf8]{inputenc}
\usepackage{libertine}   % Libertine (serif) + Biolinum (sans) -- pareado, perto do Source Serif + IBM Plex Sans do preview

% Cores emprestadas do preview (paper / paper-ink / paper-ink-muted)
\definecolor{ink}{HTML}{211F1A}
\definecolor{muted}{HTML}{7A715F}
\definecolor{hline}{HTML}{E1DCCD}
\definecolor{bodytext}{HTML}{3A362E}

\color{ink}

% Links discretos, sem azul -- igual ao preview
\hypersetup{
    colorlinks=true,
    linkcolor=ink,
    urlcolor=ink,
    citecolor=ink,
    pdfborder={0 0 0}
}

% Section formatting: sans, pequeno, com regra fina e clara (não a titlerule grossa de antes)
\titleformat{\section}
  {\normalfont\sffamily\bfseries\small\color{ink}}
  {}{0em}{}
  [\vspace{2pt}{\color{hline}\titlerule[0.6pt]}]
\titlespacing{\section}{0pt}{4pt}{2pt}

\pagestyle{empty}
\setlength{\parindent}{0pt}
\setlength{\parskip}{0.5pt}
\setlist[itemize]{itemsep=0.5pt, parsep=0.5pt, topsep=1pt, leftmargin=15pt}

% Evitar quebras de página desnecessárias
\clubpenalty=10000
\widowpenalty=10000
\displaywidowpenalty=10000
\raggedbottom

% Comandos auxiliares para replicar os blocos do preview
% \entryhead{titulo}{data}
\newcommand{\entryhead}[2]{%
  \noindent{\rmfamily\bfseries #1}\hfill{\sffamily\small\color{muted}#2}\\
}
% \entrysub{subtitulo}
\newcommand{\entrysub}[1]{%
  {\sffamily\itshape\small\color{muted}#1}\par
}

\begin{document}

% HEADER
\begin{center}
    {\rmfamily\Huge\bfseries Allen Tran}\\[2pt]
    {\sffamily\large\color{muted}Software Engineer}\\[2pt]
    {\sffamily\small\color{muted}
      allen.tran@email.com \textperiodcentered\ 
      (555) 123-4567 \textperiodcentered\ 
      \href{https://linkedin.com/in/allen-tran}{LinkedIn} \textperiodcentered\ 
      \href{https://github.com/allentran}{GitHub} \textperiodcentered\ 
      California%
    }
\end{center}

\vspace{4pt}
\noindent{\color{ink}\rule{\linewidth}{1.1pt}}
\vspace{1pt}

\section{EDUCATION}

\entryhead{Arizona State University}{2017 -- Sep 2022}
\entrysub{Bachelor of Computer Science \& Master of Computer Science}
{\sffamily\small\color{bodytext}
\begin{itemize}
    \item Focus on Artificial Intelligence
\end{itemize}
}

\section{TECHNICAL SKILLS}

{\sffamily\small\color{bodytext}
\textbf{\color{ink}Languages:} Python, Go (Golang), Java, JavaScript, TypeScript, C++, C\#, HTML/CSS \\[4pt]
\textbf{\color{ink}Technologies:} Node.js, React.js, Angular, Next.js, Express.js, AWS, Google Cloud Platform, SQL, Git, APIs, Docker, MongoDB, PostgreSQL, Heroku, Postman, REST, Prisma, Mongoose, Firebase, DynamoDB, GraphQL
}

\section{PROFESSIONAL EXPERIENCE}

\entryhead{Capital One}{June 2022 -- August 2022}
\entrysub{Software Engineer Intern \hfill Chicago, IL}
{\sffamily\small\color{bodytext}
\begin{itemize}
    \item Designed \& developed full-stack application that eliminates the need for engineers to hardcode configuration files \& reduces several minutes of deployment time per acquisition partner
    \item Deployed cloud infrastructure for Lambda functions \& S3 buckets by building pipelines that trigger upon every PR
    \item Implemented enterprise-grade SSO security layer on top of application by writing a custom library service that leverages hashed code challenges, access tokens, redirect URIs, and internal protocols
    \item Technologies used: Go (Golang), React.js, Node.js, Java, AWS Lambda | S3 | CloudFront, YAML, Jenkins
\end{itemize}
}

\entryhead{Cadent}{March 2022 -- May 2022}
\entrysub{Software Engineer Intern \hfill San Jose, CA}
{\sffamily\small\color{bodytext}
\begin{itemize}
    \item Built full-stack diagnostics dashboard that visualizes advertisement metric delivery on customer campaigns
    \item Merged code into production allowing developers company-wide to use my dashboard and save 15\% of debugging
    \item Spearheaded entire design process of system architecture and was approved by the director of engineering
    \item Graphed several million data points by building 2 new endpoints that efficiently query data from Google Cloud
    \item Technologies used: Google Cloud BigQuery, Java + Spring Boot, Angular, TypeScript, amCharts
\end{itemize}
}

\entryhead{Software \& Computer Engineering Society}{June 2021 -- August 2021}
\entrysub{Software Engineer Intern \hfill San Jose, CA}
{\sffamily\small\color{bodytext}
\begin{itemize}
    \item Replaced gRPC infrastructure with REST APIs to fix 2D \& 3D printing services used by 1,500+ club members
    \item Leveraged cloud-based queue system to read S3 hosted files, increasing PDF size support from 256KB to 10MB
    \item Implemented unit tests across frontend pages using stubs, speeding code deployment by 15\% \& reducing bug tickets
    \item Technologies used: AWS SQS \& S3, Node.js, Express.js, Sinon.js, Mocha.js, OpenVPN
\end{itemize}
}

\entryhead{Tesla}{January 2023 -- March 2023}
\entrysub{Software Engineer Intern \hfill Fremont, CA}
{\sffamily\small\color{bodytext}
\begin{itemize}
    \item Optimized latency period of data transfer through web socket by 900ms by removing unnecessary API calls
    \item Migrated health check pinging functionality to newly built REST API system and developed a new package that allows easy integration for future health checks to be implemented across the codebase
    \item Technologies used: GoLang, Typescript, React.js, Zap, Kafka, GitHub Actions
\end{itemize}
}

\section{PROJECTS}

\entryhead{Drop It | Google Drive Replica}{\href{https://github.com/rlevidev/repository}{GitHub Repository}}
{\sffamily\small\color{bodytext}
\begin{itemize}
    \item Individually developed dropbox application that allows users to upload, edit, view and download files locally
    \item Constructed Node.js API server enabling file upload to S3 bucket and user information into MySQL database
\end{itemize}
}

\section{LANGUAGES}

{\sffamily\small\color{bodytext}
\textbf{\color{ink}Portuguese:} Native \textperiodcentered\ \textbf{\color{ink}English:} Fluent (Professional Working Proficiency)
}

\section{CERTIFICATIONS}

{\sffamily\small\color{bodytext}
AWS Certified Solutions Architect Associate (2023) \textperiodcentered\ Google Cloud Professional Data Engineer (2022) \textperiodcentered\ Kubernetes Administrator Certified (CKA) (2021)
}

\end{document}
```

[template_curriculo_v2.pdf](Template%20de%20Curr%C3%ADculo/template_curriculo_v2.pdf)