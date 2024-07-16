import pandas as pd
from sklearn.preprocessing import LabelEncoder
import numpy as np
import joblib

categorical_cols = ['asn', 'country_code', 'os']
le=joblib.load('/code/app/labelEncoder.joblib')
XGBoostCF = joblib.load("/code/app/XGBModel.pkl")
def df_to_ml(df):
    old_df = pd.DataFrame(df)
    df["port_count"] = 0
    df["tcp_count"] = 0
    df["udp_count"] = 0
    for index, row in df.iterrows():
        if row["port_data"] == None:
            continue
        else:
            df["port_count"][index] = len(row["port_data"])
            for i in row["port_data"]:
                if i["protocol"] == "tcp":
                    df["tcp_count"][index] += 1
                df["udp_count"][index] = df["port_count"][index] - df["tcp_count"][index]
    

    df[categorical_cols] = df[categorical_cols].apply(lambda col: le.fit_transform(col))
    X = df.drop([ "ip","port_data"], axis=1)
    y_pred = XGBoostCF.predict(X)
    print(old_df)
    old_df["status"] = y_pred
    print(old_df)
    out = old_df.to_json(orient='records', lines=True)
    
    return out
