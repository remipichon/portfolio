import React, {useEffect, useState} from "react";


type Job = {
    namespace: string;
    name: string;
    lastSuccessfullyRunStarTime?: Date;
    lastStatus: {
        type: string;
        message?: string;
    };
    lastSuccessfullyRunCompletionTime?: Date;
};

export function App() {
    const [jobs, setJobs] = useState<Job[]>([]);
    const [loading, setLoading] = useState(true);
    const [lastFetchJobs, setLastFetchJobs] = useState<Date | null>(null);
    const [error, setError] = useState<{ code: number, message: string } | null>(null);

    useEffect(() => {
        // Fetch initially
        fetchJobs().catch(console.error);

        // Set up visibility change listener to refresh the list only when the tab is active
        const handleVisibilityChange = () => {
            if (document.visibilityState === 'visible') {
                console.info("Tab just got visible, perform reload");
                fetchJobs().catch(console.error);;
            }
        };
        console.info("Add visibilitychange event listener");
        document.addEventListener('visibilitychange', handleVisibilityChange);

        // Set up polling to refresh the list every 5s (when the tab is active)
        const intervalId = setInterval(() => {
            if (document.visibilityState === 'visible') {
                console.info("Polling fetchJobs every 5s");
                fetchJobs().catch(console.error);;
            }
        }, 5000);

        // Cleanup error notification
        if (error) {
            const timer = setTimeout(() => setError(null), 20000);
            return () => clearTimeout(timer);
        }

        // Cleanup auto refresh Job list on unmount (tab is no longer active)
        return () => {
            console.info("Disable Polling fetchJobs every 5s");
            document.removeEventListener('visibilitychange', handleVisibilityChange);
            clearInterval(intervalId);
        };
    }, [error]);

    // Refresh the list when the tab get the focus
    const handleVisibilityChange = () => {
        if (document.visibilityState === 'visible') {
            console.log("Tab just got visible, perform reload")
            fetchJobs()
        }
    };

    const fetchJobs = async () => {
        setLoading(true);
        try {
            const res = await fetch("/list");
            if (!res.ok) {
                const text = await res.text();
                throw new Error(`Error ${res.status}: ${text}`);
            }
            const jobsRaw = await res.json();
            const jobs: Job[] = jobsRaw["jobs"].map(parseJob);
            setJobs(jobs);
            setLastFetchJobs(new Date());
        } catch (err: any) {
            console.error("Fetch failed:", err);
            const match = err.message.match(/Error (\d+): (.*)/);
            if (match) {
                setError({ code: parseInt(match[1]), message: JSON.parse(match[2]).error });
            } else {
                setError({ code: 500, message: "Unknown error" });
            }
        } finally {
            setLoading(false);
        }
    };

    const parseJob = (raw: any): Job => ({
        ...raw,
        lastSuccessfullyRunStarTime: raw.lastSuccessfullyRunStarTime ? new Date(raw.lastSuccessfullyRunStarTime) : undefined,
        lastSuccessfullyRunCompletionTime: raw.lastSuccessfullyRunCompletionTime ? new Date(raw.lastSuccessfullyRunCompletionTime) : undefined,
    });

    const performJobAction = async (path: string) => {
        try {
            const res = await fetch(path);
            if (!res.ok) {
                const text = await res.text();
                throw new Error(`Error ${res.status}: ${text}`);
            }
            await fetchJobs();
        } catch (err: any) {
            console.error("Job action failed:", err);
            const match = err.message.match(/Error (\d+): (.*)/);
            if (match) {
                setError({ code: parseInt(match[1]), message: JSON.parse(match[2]).error });
            } else {
                setError({ code: 500, message: "Unknown error" });
            }
        }
    };

    const runJob = (namespace: string, name: string) =>
        performJobAction(`/run/${namespace}/${name}`);

    const killJob = (namespace: string, name: string) =>
        performJobAction(`/kill/${namespace}/${name}`);

    return (
        <div style={{padding: "2rem", fontFamily: "Arial, sans-serif"}}>
            {error && (
                <div style={errorStyle}>
                    <span>
                        ⚠ Error {error.code}: {error.message}
                    </span>
                    <button onClick={() => setError(null)} style={dismissButtonStyle}>×</button>
                </div>
            )}
            <h2>Job Assistant</h2>
            {lastFetchJobs && (
                <p>Last fetch: {lastFetchJobs.toLocaleTimeString()}</p>
            )}            {loading ? (
                <p>Loading jobs...</p>
            ) : (
                <table style={{width: "100%", borderCollapse: "collapse"}}>
                    <thead>
                    <tr>
                        <th style={thStyle}>Namespace</th>
                        <th style={thStyle}>Name</th>
                        <th style={thStyle}>Status</th>
                        <th style={thStyle}>Start time</th>
                        <th style={thStyle}>Completion time</th>
                        <th style={thStyle}>Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {jobs.map((job) => (
                        <tr key={`${job.namespace}-${job.name}`}>
                            <td style={tdStyle}>{job.namespace}</td>
                            <td style={tdStyle}>{job.name}</td>
                            <td style={tdStyle}>{job.lastStatus?.type} - {job.lastStatus?.message}</td>
                            <td style={tdStyle}>{job.lastSuccessfullyRunStarTime?.toLocaleTimeString()}</td>
                            <td style={tdStyle}>{job.lastSuccessfullyRunCompletionTime?.toLocaleTimeString()}</td>
                            <td style={tdStyle}>
                                <button
                                    onClick={() => runJob(job.namespace, job.name)}
                                    disabled={job.lastStatus.type === "Running"}
                                    style={buttonStyle}
                                >
                                    Run
                                </button>
                                <button
                                    onClick={() => killJob(job.namespace, job.name)}
                                    disabled={job.lastStatus.type !== "Running"}
                                    style={{
                                        ...buttonStyle,
                                        marginLeft: "0.5rem"
                                    }}
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

const errorStyle: React.CSSProperties = {
    backgroundColor: "#ffe0e0",
    color: "#a00",
    padding: "10px",
    border: "1px solid #f5c2c2",
    borderRadius: "4px",
    marginBottom: "1rem",
    position: "relative",
};

const dismissButtonStyle: React.CSSProperties = {
    position: "absolute",
    right: "10px",
    top: "5px",
    background: "transparent",
    border: "none",
    fontSize: "18px",
    cursor: "pointer",
};

export default App;
