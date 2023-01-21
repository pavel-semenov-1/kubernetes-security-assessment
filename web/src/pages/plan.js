import React, { useState, useEffect } from "react";
import "../style/plan.css"

const planStyles = {
    height: '100%',
    width: '100%',
}

const boxStyles = {
    borderRadius: 10,
    border: '4px solid #005618',
    padding: 15,
    height: '100%',
    width: '40%',
}

const listStyles = {
    overflow: 'auto', 
    height: '90%',
}

const Entry = ({ id, html_url, title, created_at, closed_at }) => {
    const creationDate = new Date(Date.parse(created_at)).toDateString();
    const closureDate = closed_at === null ? "???" : new Date(Date.parse(closed_at)).toDateString() ;
    return (
        <li key={id}>
            <a href={html_url}>
                {title}
            </a>
            <p>
                {creationDate} - {closureDate}
            </p>
        </li>
    )
}

const Plan = () => {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const githubApiRepoUrl = process.env.REACT_APP_GITHUB_API_REPO_URL;
    const githubToken = process.env.REACT_APP_GITHUB_TOKEN;
    const githubThesisLableId = process.env.REACT_APP_GITHUB_THESIS_LABEL_ID;

    useEffect(() => {
        fetch(`${githubApiRepoUrl}issues?state=all`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${githubToken}`
            }
        })
        .then((response) => response.json())
        .then((data) => setData(data
            .filter(issue => issue.labels
                .some(label => label.id == githubThesisLableId))))
        .catch((error) => {
            setError(true);
            setData(null);
            console.error(error);
        })
        .finally(() => {
            setLoading(false);
        });
    }, []);

    return loading ? <p>Loading...</p> : 
        error ? <p>Error fetching data</p> : (
        <div style={planStyles} className={'d-flex flex-row justify-content-around'}>
            <div style={boxStyles} className={'done'}>
                <h2 className="text-center">DONE</h2>
                <ul style={listStyles}>
                    {
                        data
                        .filter(issue => issue.state === "closed")
                        .sort((a, b) => Date.parse(b.closed_at) - Date.parse(a.closed_at))
                        .map(issue => (<Entry
                            {...issue}
                        ></Entry>))
                    }
                </ul>
            </div>
            <div style={boxStyles} className={'todo'}>
                <h2 className="text-center">TODO</h2>
                <ul style={listStyles}>
                    {
                        data
                        .filter(issue => issue.state === "open")
                        .sort((a, b) => Date.parse(b.created_at) - Date.parse(a.created_at))
                        .map(issue => (<Entry
                            {...issue}
                        ></Entry>))
                    }
                </ul>
            </div>
        </div>
    )
}

export default Plan