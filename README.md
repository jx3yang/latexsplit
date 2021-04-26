# latexsplit
Compile custom sections of a LaTeX document into different PDFs

## Problem Statement
Course instructors are choosing platforms such as [Crowdmark](https://crowdmark.com/) as a means for students to 
upload their homework, as well as a means for TAs to have tools provided by those platforms to facilitate grading.
However, Crowdmark has an absolutely terrible design for the user experience of uploading homework.

As of April 25 2021, the design of Crowdmark's upload page is as follows:
- For each homework assignment, instructors may create sections to which students can upload PDF files, JPEG, or PNG
- In the case of PDF files, students can upload a whole document to one section, then drag and drop each page of the document to a different section

For the majority of Math/Engineering students, the preferred method for writing homework assignment is most probably using LaTeX. Furthermore, 
it is probably also true that they do not create a new LaTeX file for each section. As a result, when they compile their homework, they will get 
a single PDF. This PDF is then uploaded to one section, and pages that belong to other sections are dragged and dropped appropriately. As one can 
imagine, this is very bothersome whenever some combination of the following is true:
- The PDF has a lot of pages
- There are many sections

In many cases, it is also necessary to perform drag and dropping while scrolling down the page. There are some life hacks to work around this problem,
such as zooming out of the page until all sections fit on the visible screen, but it is still quite difficult if your screen is small (e.g. laptop).
Other methods involve uploading the same PDF to each section, then removing the irrelevant pages depending on the section. This eliminates the need 
for painful drag and dropping (+ scrolling), but introduces the need to remove pages. Moreover, it does not work as well if your PDF has a lot of pages.
Most importantly, those methods are most useful for the initial submission; if you need to modify your submission, then you still need to 
upload a whole PDF and do some painful page removal.

## Proposed Solution
A tool that will automatically split the output PDF from the LaTeX compiler into appropriate sections. Suppose for example that we want to 
split our PDF into two sections, then instead of simply outputing `<doc_name>.pdf`, we will have 
- `<doc_name>.pdf`
- `1_<doc_name>.pdf`
- `2_<doc_name>.pdf`

where `<i>_<doc_name>.pdf` represents the i-th section resulting from the split.

To tell the tool where exactly to split, inside the LaTeX source code, write a comment 
```latex
\begin{...}
...
\end{...}

% latexsplit
\newpage

\begin{...}
...
\end{...}
```

where the split will occur between the page where the line `% latexsplit` is located and the following page.

### Example
```latex
section 1
\newpage
section 1
\newpage
section 1
% latexsplit
\newpage

section 2
\newpage
section 2
% latexsplit
section 2
\newpage

section 3
```

In the above example, we would have 3 sections: 

- first section has 3 pages (pages 1-3)
- second section has 2 pages (pages 4-5), note that even though `% latexsplit` is in the middle of a page, the split will 
occur between that page and the subsequent page
- third section has 1 page (page 6)
