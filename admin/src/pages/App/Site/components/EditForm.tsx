import React, { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-form';
import { ProFormDigit } from '@ant-design/pro-form';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressSiteInfo,
  CtlIngressSiteSave,
  CtlIngressNamespaceMap,
  CtlIngressActionBackendMap,
  CtlIngressActionStaticMap,
  DataMapAppDataIngressActionType,
  DataMapAppDataIngressSiteAuthNeed,
  DataMapRuleStringOp,
} from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Spin } from 'antd';
import { Access } from 'umi';

export type EditFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  id?: string;
};

const EditForm: React.FC<EditFormProps> = (props) => {
  const [xIsLoading, setIsLoading] = useState<boolean>(true);

  const xFormRef = useRef<ProFormInstance<API.AppDataIngressSite>>();
  const [xNeedAuth, setNeedAuth] = useState<boolean>(false);
  const [xIngressActionValueMapData, setIngressActionValueMapData] = useState<any>();

  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const loadNamespaceMapData = async () => {
    const dataList: any[] = [];

    const xReq = {};

    const xRespData = await CtlIngressNamespaceMap(xReq);

    if (xRespData && xRespData.data) {
      xRespData.data.forEach((xItem) => {
        dataList.push({ label: xItem.data_name, value: xItem.data_id });
      });
    }

    return dataList;
  };

  const loadIngressActionValueMapData = async () => {
    const xDataList: any[] = [];

    const xActionType = xFormRef.current?.getFieldValue('action_type');

    const xNamespaceId = xFormRef.current?.getFieldValue('namespace_id');
    const xReq = { data_id: xNamespaceId };

    if (!xNamespaceId) {
      return;
    }

    if (xActionType == 'backend') {
      const xRespData = await CtlIngressActionBackendMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    if (xActionType == 'static') {
      const xRespData = await CtlIngressActionStaticMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    setIngressActionValueMapData(xDataList);
  };

  const loadAuthFields = () => {
    const xNeedAuthFlag = xFormRef.current?.getFieldValue('auth_need') == 'yes';
    setNeedAuth(xNeedAuthFlag);
  };

  const loadDataInfo = async (id?: string) => {
    setIsLoading(true);

    if (id) {
      const respData = await CtlIngressSiteInfo({ data_id: id });
      if (respData && respData.error_code === '0') {
        if (respData.data) {
          const xDataInfo = respData.data;
          xFormRef.current?.setFieldsValue(xDataInfo);
          loadAuthFields();
          await loadIngressActionValueMapData();
        }
      }
    }

    setIsLoading(false);
  };

  // ??????
  useEffect(() => {
    if (props.modalVisible) {
      loadDataInfo(props.id);
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.id, props.modalVisible]);

  const onSaveData = async (xData: API.AppDataIngressSite) => {
    xData.data_id = props.id;

    const xResp = await CtlIngressSiteSave(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressSite>
      title={'????????????'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressSite);
      }}
      formRef={xFormRef}
    >
      <Spin spinning={xIsLoading} delay={100}>
        <ProFormSelect
          name="namespace_id"
          label="????????????"
          placeholder="?????????????????????"
          rules={[{ required: true, message: '?????????????????????.' }]}
          request={loadNamespaceMapData}
          fieldProps={{
            showSearch: true,
            onChange: (value) => {
              xFormRef.current?.resetFields();
              xFormRef.current?.setFieldsValue({ namespace_id: value });
            },
          }}
        />

        <ProFormText
          name="title"
          label="????????????"
          placeholder="?????????????????????"
          rules={[{ required: true, message: '?????????????????????.' }]}
        />

        <ProFormDigit
          name="order_no"
          label="????????????"
          placeholder="?????????????????????"
          rules={[{ required: true, message: '?????????????????????.' }]}
        />

        <ProFormSelect
          name="match_op"
          label="??????????????????"
          valueEnum={DataMapRuleStringOp}
          placeholder="???????????????????????????"
          rules={[{ required: true, message: '???????????????????????????.' }]}
        />

        <ProFormText
          name="match_value"
          label="??????????????????"
          placeholder="???????????????????????????"
          rules={[{ required: true, message: '???????????????????????????.' }]}
        />

        <ProFormSelect
          name="auth_need"
          label="??????????????????"
          valueEnum={DataMapAppDataIngressSiteAuthNeed}
          placeholder="?????????????????????????????????"
          rules={[{ required: true, message: '?????????????????????????????????.' }]}
          fieldProps={{
            onChange() {
              loadAuthFields();
            },
          }}
        />

        <Access key={'auth_need_data'} accessible={xNeedAuth}>
          <ProFormText
            name="auth_user_name"
            label="?????????????????????"
            placeholder="??????????????????????????????"
            rules={[{ required: false, message: '??????????????????????????????.' }]}
          />

          <ProFormText
            name="auth_password"
            label="??????????????????"
            placeholder="???????????????????????????"
            rules={[{ required: false, message: '???????????????????????????.' }]}
          />
        </Access>

        <ProFormSelect
          name="action_type"
          label="????????????"
          valueEnum={DataMapAppDataIngressActionType}
          placeholder="?????????????????????"
          rules={[{ required: true, message: '?????????????????????.' }]}
          fieldProps={{
            onChange() {
              const xFormDataVal = xFormRef.current?.getFieldsValue();
              if (xFormDataVal) {
                xFormDataVal.action_value = undefined;
                xFormRef.current?.setFieldsValue(xFormDataVal);
              }
              loadIngressActionValueMapData();
            },
          }}
        />

        <ProFormSelect
          name="action_value"
          label="????????????"
          placeholder="?????????????????????"
          fieldProps={{ options: xIngressActionValueMapData, showSearch: true }}
          rules={[{ required: true, message: '?????????????????????.' }]}
        />
        <ProFormSelect
          name="state"
          label="??????"
          valueEnum={{
            enable: { text: '??????', status: 'enable' },
            disable: { text: '??????', status: 'disable' },
          }}
          placeholder="???????????????"
          rules={[{ required: true, message: '???????????????.' }]}
        />
      </Spin>
    </ModalForm>
  );
};

export default EditForm;
