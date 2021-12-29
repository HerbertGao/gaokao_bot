package com.herbertgao.telegram.database.service;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.herbertgao.telegram.database.entity.UserTemplate;
import com.herbertgao.telegram.database.mapper.UserTemplateMapper;
import com.herbertgao.telegram.util.IdUtil;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * 用户模板服务
 *
 * @author HerbertGao
 * @date 2020-04-08
 */
@Service
public class UserTemplateService extends ServiceImpl<UserTemplateMapper, UserTemplate> {

    /**
     * 通过用户ID获取模板列表
     *
     * @param userId 用户id
     * @return {@link List}<{@link UserTemplate}>
     */
    public List<UserTemplate> getUserTemplateListByUserId(Long userId) {
        QueryWrapper<UserTemplate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(UserTemplate::getUserId, userId);
        return this.list(wrapper);
    }

    /**
     * 插入默认模板和用户名
     *
     * @param userId       用户id
     * @param examUsername 用户名
     * @return {@link Long}
     */
    public Long insertDefaultTemplateWithUsername(Long userId, String examUsername) {
        UserTemplate defaultTemplate = getDefaultTemplate();
        UserTemplate userTemplate = new UserTemplate();
        Long id = IdUtil.getId();
        userTemplate.setId(id);
        userTemplate.setUserId(userId);
        userTemplate.setTemplateName("@" + examUsername);
        userTemplate.setTemplateContent(defaultTemplate.getTemplateContent() + "\n@" + examUsername);
        this.save(userTemplate);
        return id;
    }

    /**
     * 插入模板和用户名
     *
     * @param userId       用户id
     * @param template     模板
     * @param templateName 模板名称
     * @return {@link Long}
     */
    public Long insertTemplateWithUsername(Long userId, String template, String templateName) {
        UserTemplate userTemplate = new UserTemplate();
        Long id = IdUtil.getId();
        userTemplate.setId(id);
        userTemplate.setUserId(userId);
        userTemplate.setTemplateName(templateName);
        userTemplate.setTemplateContent(template);
        this.save(userTemplate);
        return id;
    }

    /**
     * 获取默认模板列表
     *
     * @return {@link UserTemplate}
     */
    public UserTemplate getDefaultTemplate() {
        return getUserTemplateListByUserId(0L).get(0);
    }

    /**
     * 通过id获取模板
     *
     * @param userId 用户id
     * @param id     id
     * @return {@link UserTemplate}
     */
    public UserTemplate getTemplateById(Long userId, String id) {
        QueryWrapper<UserTemplate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(UserTemplate::getId, Long.parseLong(id))
                .eq(UserTemplate::getUserId, userId);
        return this.getOne(wrapper);
    }

    /**
     * 删除
     *
     * @param userId 用户id
     * @param id     id
     */
    public void remove(Long userId, String id) {
        QueryWrapper<UserTemplate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(UserTemplate::getId, Long.parseLong(id))
                .eq(UserTemplate::getUserId, userId);
        this.remove(wrapper);
    }

    /**
     * 更新
     *
     * @param template 模板
     */
    public void update(UserTemplate template) {
        this.updateById(template);
    }
}
