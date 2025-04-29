import React, { useEffect, useState } from "react";

type Job = {
    name: string;
    namespace: string;
    status: "suspended" | "running" | "failed" | "scheduled";
};

export function App() {
    const [jobs, setJobs] = useState<Job[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchJobs = async () => {
        setLoading(true);
        const res = await fetch("/list");
        const data = await res.json();
        setJobs(
            data["jobs"].map((job: any) => ({
                name: job.name,
                namespace: job.namespace,
                status: getJobStatus(job),
            }))
        );
        setLoading(false);
    };

    useEffect(() => {
        fetchJobs();
    }, []);

    const getJobStatus = (job: any): Job["status"] => {
        if (job.spec?.suspend === true) return "suspended";
        if (job.status?.active > 0) return "running";
        if (job.status?.failed > 0) return "failed";
        return "scheduled";
    };

    const runJob = async (namespace: string, name: string) => {
        await fetch(`/run/${namespace}/${name}`);
        fetchJobs();
    };

    const killJob = async (namespace: string, name: string) => {
        await fetch(`/kill/${namespace}/${name}`);
        fetchJobs();
    };

    return (
        <div style={{ padding: "2rem", fontFamily: "Arial, sans-serif" }}>
            <h2>Job Assistant</h2>
            {loading ? (
                <p>Loading jobs...</p>
            ) : (
                <table style={{ width: "100%", borderCollapse: "collapse" }}>
                    <thead>
                    <tr>
                        <th style={thStyle}>Namespace</th>
                        <th style={thStyle}>Name</th>
                        <th style={thStyle}>Status</th>
                        <th style={thStyle}>Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {jobs.map((job) => (
                        <tr key={`${job.namespace}-${job.name}`}>
                            <td style={tdStyle}>{job.namespace}</td>
                            <td style={tdStyle}>{job.name}</td>
                            <td style={tdStyle}>{job.status}</td>
                            <td style={tdStyle}>
                                <button
                                    onClick={() => runJob(job.namespace, job.name)}
                                    disabled={job.status === "running" || job.status === "scheduled"}
                                    style={buttonStyle}
                                >
                                    Run
                                </button>
                                <button
                                    onClick={() => killJob(job.namespace, job.name)}
                                    disabled={job.status !== "running"}
                                    style={{ ...buttonStyle, marginLeft: "0.5rem" }}
                                >
                                    Kill
                                </button>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
        </div>
    );
}

const thStyle: React.CSSProperties = {
    textAlign: "left",
    padding: "8px",
    backgroundColor: "#f2f2f2",
    borderBottom: "1px solid #ccc",
};

const tdStyle: React.CSSProperties = {
    padding: "8px",
    borderBottom: "1px solid #eee",
};

const buttonStyle: React.CSSProperties = {
    padding: "6px 12px",
    fontSize: "14px",
    borderRadius: "4px",
    border: "1px solid #ccc",
    backgroundColor: "#f5f5f5",
    cursor: "pointer",
};

export default App;
