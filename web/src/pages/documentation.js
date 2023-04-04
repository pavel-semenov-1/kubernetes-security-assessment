import React from "react";

const documentationStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
}

const Documentation = () => {
    const documentList = process.env.DOCUMENT_LIST;
    let documents = [];
    if (documentList != undefined) {
        documents = documentList.trim().split("~").map(document => {
            const [title, href] = document.split(":")
            return <li><a href={href} target="_blank">{title}</a></li>
        })
    }
    
    return (
        <div style={documentationStyles} className={'documentation'}>
            <h2 className="text-center">Documents</h2>
            <ul>
                {documents}
            </ul>
        </div>
    )
}

export default Documentation