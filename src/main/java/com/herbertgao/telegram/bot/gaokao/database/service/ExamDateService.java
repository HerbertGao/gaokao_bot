package com.herbertgao.telegram.bot.gaokao.database.service;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.herbertgao.telegram.bot.gaokao.database.domain.ExamDate;
import com.herbertgao.telegram.bot.gaokao.database.mapper.ExamDateMapper;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

/**
 * 考试日期服务
 *
 * @author HerbertGao
 * @date 2019-06-09
 */
@Service
public class ExamDateService extends ServiceImpl<ExamDateMapper, ExamDate> {

    /**
     * 得到考试列表
     *
     * @param now  现在
     * @param desc 倒序
     * @return {@link List}<{@link ExamDate}>
     */
    public List<ExamDate> getExamList(LocalDateTime now, boolean desc) {
        QueryWrapper<ExamDate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .le(ExamDate::getExamYearBeginDate, now)
                .ge(ExamDate::getExamYearEndDate, now)
                .eq(ExamDate::getIsDelete, false);
        if (desc) {
            wrapper.lambda().orderByDesc(ExamDate::getExamYearBeginDate);
        }
        return this.list(wrapper);
    }

    /**
     * 被一年考试
     *
     * @param year 年份
     * @return {@link List}<{@link ExamDate}>
     */
    public List<ExamDate> getExamByYear(Integer year) {
        QueryWrapper<ExamDate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(ExamDate::getExamYear, year)
                .eq(ExamDate::getIsDelete, false);
        return this.list(wrapper);
    }
}
