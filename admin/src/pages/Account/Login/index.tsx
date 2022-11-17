import { LockOutlined, SafetyCertificateOutlined, UserOutlined } from '@ant-design/icons';
import React, { useState } from 'react';
import { ProFormText, LoginForm } from '@ant-design/pro-form';
import { useIntl, history, FormattedMessage, SelectLang, useModel } from 'umi';
import Footer from '@/components/Footer';

import styles from './index.less';
import {
  CtlAccountLogin,
  GetCheckCodeImageId,
  setLoginAccessToken,
} from '@/services/ant-design-pro/api';
import CheckCodeImage from '@/components/CheckCodeImage';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';

const Login: React.FC = () => {
  const [submitting, setSubmitting] = useState(false);
  const intl = useIntl();
  const { initialState, setInitialState } = useModel('@@initialState');

  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s: any) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };

  const handleSubmit = async (values: API.AccountLoginRequest) => {
    try {
      setSubmitting(true);

      values.image_code_id = GetCheckCodeImageId();

      const respData = await CtlAccountLogin({ ...values });

      setSubmitting(false);

      if (respData.error_code === '0') {
        ShowSuccessMessage(respData.error_message);

        //设置登录信息
        setLoginAccessToken(respData.data);
        await fetchUserInfo();

        if (!history) {
          return;
        }

        const { query } = history.location;

        const { redirect } = query as { redirect: string };

        history.push(redirect || '/');

        return;
      }

      ShowErrorMessage(respData.error_message);
    } catch (error) {
      ShowErrorMessage('pages.login.failure');
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.lang} data-lang>
        {SelectLang && <SelectLang />}
      </div>
      <div className={styles.content}>
        <LoginForm
          logo={<img alt="logo" src="/logo.svg" />}
          title={intl.formatMessage({ id: 'pages.layouts.userLayout.sitename' })}
          subTitle={intl.formatMessage({ id: 'pages.layouts.userLayout.title' })}
          submitter={{
            searchConfig: {
              submitText: intl.formatMessage({
                id: 'pages.login.submit.login',
                defaultMessage: '登录',
              }),
            },
            render: (_, dom) => dom.pop(),
            submitButtonProps: {
              loading: submitting,
              size: 'large',
              style: {
                width: '100%',
              },
            },
          }}
          onFinish={async (values) => {
            await handleSubmit(values as API.AccountLoginRequest);
          }}
        >
          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined className={styles.prefixIcon} />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.login.username.placeholder',
              defaultMessage: '用户名',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.login.username.required"
                    defaultMessage="请输入用户名!"
                  />
                ),
              },
            ]}
          />

          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className={styles.prefixIcon} />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.login.password.placeholder',
              defaultMessage: '密码',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.login.password.required"
                    defaultMessage="请输入密码！"
                  />
                ),
              },
            ]}
          />

          <ProFormText
            name="image_code_data"
            fieldProps={{
              size: 'large',
              prefix: <SafetyCertificateOutlined />,
              addonAfter: <CheckCodeImage />,
            }}
            placeholder={intl.formatMessage({ id: 'pages.input.checkimagecode.placeholder' })}
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.input.checkimagecode.required" />,
              },
            ]}
          />
        </LoginForm>
      </div>
      <Footer />
    </div>
  );
};

export default Login;
